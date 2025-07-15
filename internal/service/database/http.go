package database

import (
	"context"
	"encoding/json"
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

// Config demonstrates the configuration registry pattern
type Config struct {
	DatabaseURL    string        `envconfig:"DATABASE_URL" default:"mock://localhost:5432/testdb"`
	MaxConnections int           `envconfig:"DATABASE_MAX_CONNECTIONS" default:"10"`
	ConnectTimeout time.Duration `envconfig:"DATABASE_CONNECT_TIMEOUT" default:"5s"`
	EnableMetrics  bool          `envconfig:"DATABASE_ENABLE_METRICS" default:"true"`
}

// Load implements the Configurable interface for the configuration registry
func (c *Config) Load() error {
	if err := envconfig.Process("", c); err != nil {
		return fmt.Errorf("failed to load database config: %w", err)
	}

	// Custom validation
	if c.MaxConnections < 1 || c.MaxConnections > 100 {
		return fmt.Errorf("DATABASE_MAX_CONNECTIONS must be between 1 and 100, got %d", c.MaxConnections)
	}

	if c.ConnectTimeout < time.Second {
		return fmt.Errorf("DATABASE_CONNECT_TIMEOUT must be at least 1 second, got %v", c.ConnectTimeout)
	}

	return nil
}

// DatabaseService demonstrates lifecycle management with external resources
type DatabaseService struct {
	config     *Config
	connection *MockConnection
	metrics    map[string]int
	mu         sync.RWMutex
}

func NewService() *DatabaseService {
	cfg := &Config{}

	// Register configuration with the global registry
	config.Register(cfg)

	return &DatabaseService{
		config:  cfg,
		metrics: make(map[string]int),
	}
}

// Name implements the API interface
func (d *DatabaseService) Name() string {
	return "database"
}

// Register implements the API interface
func (d *DatabaseService) Register(app server.Server) error {
	router := app.Router().(*chi.Mux)

	// Mount database routes
	router.Route("/api/database", func(r chi.Router) {
		r.Get("/users", d.getUsers)
		r.Post("/users", d.createUser)
		r.Get("/health", d.healthCheck)
		r.Get("/metrics", d.getMetrics)
	})

	return nil
}

// Start implements the LifecycleAPI interface
func (d *DatabaseService) Start(ctx context.Context) error {
	log.Info().
		Str("url", d.config.DatabaseURL).
		Int("max_connections", d.config.MaxConnections).
		Dur("timeout", d.config.ConnectTimeout).
		Msg("Connecting to database")

	// Validate config was loaded properly
	if d.config.ConnectTimeout == 0 {
		d.config.ConnectTimeout = 5 * time.Second
		log.Warn().Msg("Using default connect timeout as config was not loaded properly")
	}

	// Create connection with timeout
	connectCtx, cancel := context.WithTimeout(ctx, d.config.ConnectTimeout)
	defer cancel()

	connection, err := NewMockConnection(connectCtx, d.config.DatabaseURL, d.config.MaxConnections)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test connection
	if err := connection.Ping(connectCtx); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	d.connection = connection
	log.Info().Msg("Database connection established successfully")
	return nil
}

// Stop implements the LifecycleAPI interface
func (d *DatabaseService) Stop(ctx context.Context) error {
	log.Info().Msg("Shutting down database service")

	if d.connection != nil {
		if err := d.connection.Close(ctx); err != nil {
			log.Error().Err(err).Msg("Error closing database connection")
			return err
		}
	}

	log.Info().Msg("Database service shutdown complete")
	return nil
}

// API Handlers demonstrate error context preservation

func (d *DatabaseService) getUsers(w http.ResponseWriter, r *http.Request) {
	if d.connection == nil {
		err := request.NewErrServerWithContext(r)
		render.Render(w, r, err)
		return
	}

	users, err := d.connection.GetUsers(r.Context())
	if err != nil {
		// Use context-aware error handling for better tracing
		errResp := request.ErrInvalidRequestWithContext(r, err)
		render.Render(w, r, errResp)
		return
	}

	d.incrementMetric("users_fetched")
	render.JSON(w, r, map[string]interface{}{
		"users": users,
		"count": len(users),
	})
}

func (d *DatabaseService) createUser(w http.ResponseWriter, r *http.Request) {
	if d.connection == nil {
		err := request.NewErrServerWithContext(r)
		render.Render(w, r, err)
		return
	}

	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errResp := request.ErrInvalidRequestWithContext(r, err)
		render.Render(w, r, errResp)
		return
	}

	if req.Name == "" || req.Email == "" {
		errResp := request.ErrInvalidRequestWithContext(r, fmt.Errorf("name and email are required"))
		render.Render(w, r, errResp)
		return
	}

	user, err := d.connection.CreateUser(r.Context(), req.Name, req.Email)
	if err != nil {
		errResp := request.ErrInvalidRequestWithContext(r, err)
		render.Render(w, r, errResp)
		return
	}

	d.incrementMetric("users_created")
	render.JSON(w, r, map[string]interface{}{
		"user":    user,
		"message": "User created successfully",
	})
}

func (d *DatabaseService) healthCheck(w http.ResponseWriter, r *http.Request) {
	if d.connection == nil {
		err := request.NewErrServerWithContext(r)
		render.Render(w, r, err)
		return
	}

	if err := d.connection.Ping(r.Context()); err != nil {
		errResp := request.ErrInvalidRequestWithContext(r, err)
		render.Render(w, r, errResp)
		return
	}

	render.JSON(w, r, map[string]interface{}{
		"status":     "healthy",
		"connection": "active",
		"timestamp":  time.Now(),
	})
}

func (d *DatabaseService) getMetrics(w http.ResponseWriter, r *http.Request) {
	if !d.config.EnableMetrics {
		err := request.NewErrNotFoundWithContext(r)
		render.Render(w, r, err)
		return
	}

	d.mu.RLock()
	metrics := make(map[string]int)
	for k, v := range d.metrics {
		metrics[k] = v
	}
	d.mu.RUnlock()

	render.JSON(w, r, map[string]interface{}{
		"metrics": metrics,
		"config": map[string]interface{}{
			"max_connections": d.config.MaxConnections,
			"database_url":    "[REDACTED]", // Don't expose sensitive data
		},
	})
}

func (d *DatabaseService) incrementMetric(name string) {
	if !d.config.EnableMetrics {
		return
	}

	d.mu.Lock()
	d.metrics[name]++
	d.mu.Unlock()
}

// MockConnection simulates a database connection for demonstration
type MockConnection struct {
	url         string
	maxConns    int
	users       []User
	isConnected bool
}

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func NewMockConnection(ctx context.Context, url string, maxConns int) (*MockConnection, error) {
	// Simulate minimal connection delay
	select {
	case <-time.After(10 * time.Millisecond):
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	// Simulate connection failure for certain URLs
	if url == "mock://fail" {
		return nil, fmt.Errorf("connection failed")
	}

	return &MockConnection{
		url:         url,
		maxConns:    maxConns,
		isConnected: true,
		users: []User{
			{ID: 1, Name: "Alice", Email: "alice@example.com"},
			{ID: 2, Name: "Bob", Email: "bob@example.com"},
		},
	}, nil
}

func (c *MockConnection) Ping(ctx context.Context) error {
	if !c.isConnected {
		return fmt.Errorf("connection closed")
	}
	return nil
}

func (c *MockConnection) GetUsers(ctx context.Context) ([]User, error) {
	if !c.isConnected {
		return nil, fmt.Errorf("connection closed")
	}
	return c.users, nil
}

func (c *MockConnection) CreateUser(ctx context.Context, name, email string) (User, error) {
	if !c.isConnected {
		return User{}, fmt.Errorf("connection closed")
	}

	user := User{
		ID:    len(c.users) + 1,
		Name:  name,
		Email: email,
	}
	c.users = append(c.users, user)
	return user, nil
}

func (c *MockConnection) Close(ctx context.Context) error {
	c.isConnected = false
	log.Info().Msg("Mock database connection closed")
	return nil
}
