package worker

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"

	"github.com/go-obvious/server"
	"github.com/go-obvious/server/config"
	"github.com/go-obvious/server/request"
)

// Config demonstrates configuration registry with worker settings
type Config struct {
	WorkerInterval    time.Duration `envconfig:"WORKER_INTERVAL" default:"30s"`
	MaxJobs          int           `envconfig:"WORKER_MAX_JOBS" default:"100"`
	EnableProcessing bool          `envconfig:"WORKER_ENABLE_PROCESSING" default:"true"`
}

// Load implements the Configurable interface
func (c *Config) Load() error {
	if err := envconfig.Process("", c); err != nil {
		return fmt.Errorf("failed to load worker config: %w", err)
	}
	
	// Custom validation
	if c.WorkerInterval < time.Second {
		return fmt.Errorf("WORKER_INTERVAL must be at least 1 second, got %v", c.WorkerInterval)
	}
	
	if c.MaxJobs < 1 || c.MaxJobs > 1000 {
		return fmt.Errorf("WORKER_MAX_JOBS must be between 1 and 1000, got %d", c.MaxJobs)
	}
	
	return nil
}

// WorkerService demonstrates background task management with lifecycle
type WorkerService struct {
	config    *Config
	stopCh    chan struct{}
	wg        sync.WaitGroup
	jobs      []Job
	processed int
	mu        sync.RWMutex
	isRunning bool
}

type Job struct {
	ID        int       `json:"id"`
	Type      string    `json:"type"`
	Data      string    `json:"data"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewService() *WorkerService {
	cfg := &Config{}
	
	// Register configuration with the global registry
	config.Register(cfg)
	
	return &WorkerService{
		config: cfg,
		stopCh: make(chan struct{}),
		jobs:   make([]Job, 0),
	}
}

// Name implements the API interface
func (w *WorkerService) Name() string {
	return "worker"
}

// Register implements the API interface
func (w *WorkerService) Register(app server.Server) error {
	router := app.Router().(*chi.Mux)
	
	// Mount worker routes
	router.Route("/api/worker", func(r chi.Router) {
		r.Get("/status", w.getStatus)
		r.Get("/jobs", w.getJobs)
		r.Post("/jobs", w.createJob)
		r.Get("/health", w.healthCheck)
	})
	
	return nil
}

// Start implements the LifecycleAPI interface
func (w *WorkerService) Start(ctx context.Context) error {
	log.Info().
		Dur("interval", w.config.WorkerInterval).
		Int("max_jobs", w.config.MaxJobs).
		Bool("enabled", w.config.EnableProcessing).
		Msg("Starting background worker")
	
	if !w.config.EnableProcessing {
		log.Info().Msg("Worker processing disabled by configuration")
		return nil
	}
	
	w.isRunning = true
	
	// Start background worker goroutine
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		w.runWorker()
	}()
	
	log.Info().Msg("Background worker started successfully")
	return nil
}

// Stop implements the LifecycleAPI interface
func (w *WorkerService) Stop(ctx context.Context) error {
	log.Info().Msg("Shutting down background worker")
	
	if !w.isRunning {
		log.Info().Msg("Worker was not running")
		return nil
	}
	
	// Signal worker to stop
	close(w.stopCh)
	w.isRunning = false
	
	// Wait for worker to finish with context timeout
	done := make(chan struct{})
	go func() {
		w.wg.Wait()
		close(done)
	}()
	
	select {
	case <-done:
		log.Info().Msg("Background worker stopped gracefully")
	case <-ctx.Done():
		log.Warn().Msg("Background worker shutdown timed out")
		return ctx.Err()
	}
	
	return nil
}

// runWorker is the main worker loop
func (w *WorkerService) runWorker() {
	ticker := time.NewTicker(w.config.WorkerInterval)
	defer ticker.Stop()
	
	log.Info().Msg("Worker loop started")
	
	for {
		select {
		case <-ticker.C:
			w.processJobs()
		case <-w.stopCh:
			log.Info().Msg("Worker loop stopping")
			return
		}
	}
}

// processJobs simulates processing pending jobs
func (w *WorkerService) processJobs() {
	w.mu.Lock()
	defer w.mu.Unlock()
	
	pendingJobs := 0
	for i := range w.jobs {
		if w.jobs[i].Status == "pending" {
			// Simulate job processing
			w.jobs[i].Status = "completed"
			w.jobs[i].UpdatedAt = time.Now()
			w.processed++
			pendingJobs++
			
			log.Debug().
				Int("job_id", w.jobs[i].ID).
				Str("type", w.jobs[i].Type).
				Msg("Job processed")
		}
	}
	
	if pendingJobs > 0 {
		log.Info().
			Int("processed", pendingJobs).
			Int("total_processed", w.processed).
			Msg("Jobs processed")
	}
}

// API Handlers

func (w *WorkerService) getStatus(w2 http.ResponseWriter, r *http.Request) {
	w.mu.RLock()
	totalJobs := len(w.jobs)
	pendingJobs := 0
	completedJobs := 0
	
	for _, job := range w.jobs {
		switch job.Status {
		case "pending":
			pendingJobs++
		case "completed":
			completedJobs++
		}
	}
	w.mu.RUnlock()
	
	render.JSON(w2, r, map[string]interface{}{
		"status": map[string]interface{}{
			"running":   w.isRunning,
			"enabled":   w.config.EnableProcessing,
			"interval":  w.config.WorkerInterval.String(),
			"max_jobs":  w.config.MaxJobs,
		},
		"metrics": map[string]interface{}{
			"total_jobs":     totalJobs,
			"pending_jobs":   pendingJobs,
			"completed_jobs": completedJobs,
			"processed":      w.processed,
		},
	})
}

func (w *WorkerService) getJobs(w2 http.ResponseWriter, r *http.Request) {
	w.mu.RLock()
	jobs := make([]Job, len(w.jobs))
	copy(jobs, w.jobs)
	w.mu.RUnlock()
	
	render.JSON(w2, r, map[string]interface{}{
		"jobs":  jobs,
		"count": len(jobs),
	})
}

func (w *WorkerService) createJob(w2 http.ResponseWriter, r *http.Request) {
	var req struct {
		Type string `json:"type"`
		Data string `json:"data"`
	}
	
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		errResp := request.ErrInvalidRequestWithContext(r, err)
		render.Render(w2, r, errResp)
		return
	}
	
	if req.Type == "" {
		errResp := request.ErrInvalidRequestWithContext(r, fmt.Errorf("job type is required"))
		render.Render(w2, r, errResp)
		return
	}
	
	w.mu.Lock()
	
	// Check job limit
	if len(w.jobs) >= w.config.MaxJobs {
		w.mu.Unlock()
		errResp := request.ErrInvalidRequestWithContext(r, fmt.Errorf("job queue is full (max: %d)", w.config.MaxJobs))
		render.Render(w2, r, errResp)
		return
	}
	
	job := Job{
		ID:        len(w.jobs) + 1,
		Type:      req.Type,
		Data:      req.Data,
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	w.jobs = append(w.jobs, job)
	w.mu.Unlock()
	
	log.Info().
		Int("job_id", job.ID).
		Str("type", job.Type).
		Msg("Job created")
	
	render.JSON(w2, r, map[string]interface{}{
		"job":     job,
		"message": "Job created successfully",
	})
}

func (w *WorkerService) healthCheck(w2 http.ResponseWriter, r *http.Request) {
	status := "healthy"
	if !w.isRunning && w.config.EnableProcessing {
		status = "unhealthy"
	}
	
	render.JSON(w2, r, map[string]interface{}{
		"status":     status,
		"running":    w.isRunning,
		"enabled":    w.config.EnableProcessing,
		"timestamp":  time.Now(),
	})
}