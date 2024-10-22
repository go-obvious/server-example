package handlers

import (
	"net/http"

	"github.com/go-chi/chi"

	"github.com/go-obvious/server"
	"github.com/go-obvious/server/api"
	"github.com/go-obvious/server/request"
)

type Ping struct {
	api.Service
}

func NewPingService(base string) *Ping {
	a := &Ping{}
	a.Service.APIName = "ping"
	a.Service.Mounts = map[string]*chi.Mux{}
	a.Service.Mounts[base] = a.Routes()
	return a
}

func (a *Ping) Register(app server.Server) error {
	if err := a.Service.Register(app); err != nil {
		return err
	}
	return nil
}

func (a *Ping) Routes() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/", a.Handler)
	return r
}

type Response struct {
	Message string `json:"message"`
}

func (a *Ping) Handler(w http.ResponseWriter, r *http.Request) {
	request.Reply(r, w, &Response{Message: "PING"}, http.StatusOK)
}
