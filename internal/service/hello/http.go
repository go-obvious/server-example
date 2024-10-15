package hello

import (
	"net/http"

	"github.com/go-chi/chi"

	"github.com/go-obvious/server"
	"github.com/go-obvious/server/api"
	"github.com/go-obvious/server/request"
)

type API struct {
	api.Service
}

func NewService(base string) *API {
	a := &API{}
	a.Service.APIName = "hello"
	a.Service.Mounts = map[string]*chi.Mux{}
	a.Service.Mounts[base] = a.Routes()
	return a
}

func (a *API) Register(app server.Server) error {
	if err := a.Service.Register(app); err != nil {
		return err
	}
	return nil
}

func (a *API) Routes() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/", a.Handler)
	return r
}

type Response struct {
	Message string `json:"message"`
}

func (a *API) Handler(w http.ResponseWriter, r *http.Request) {
	request.Reply(r, w, &Response{Message: "hello"}, http.StatusOK)
}
