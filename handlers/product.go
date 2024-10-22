package handlers

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi"

	"github.com/go-obvious/server"
	"github.com/go-obvious/server-example/types"
	"github.com/go-obvious/server/api"
	"github.com/go-obvious/server/request"
)

type Products struct {
	api.Service
	store types.Store
}

func NewProductService(path string, s types.Store) *Products {
	a := &Products{
		Service: api.Service{
			APIName: "products",
			Mounts:  map[string]*chi.Mux{},
			Router:  nil,
		},
		store: s,
	}
	a.Service.Mounts[path] = a.Routes()
	return a
}

func (a *Products) Register(app server.Server) error {
	if err := a.Service.Register(app); err != nil {
		return err
	}
	return nil
}

func (a *Products) Routes() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/", a.AllProducts)
	r.Get("/{id}", a.GetProduct)
	r.Put("/{id}", a.PutProduct)
	r.Delete("/{id}", a.DeleteProduct)
	return r
}

func (a *Products) AllProducts(w http.ResponseWriter, r *http.Request) {
	next := request.QS(r, "next")
	var cursor *string
	if strings.TrimSpace(next) == "" {
		cursor = &next
	}

	productRange, err := a.store.All(r.Context(), cursor)
	if err != nil {
		request.Reply(r, w, err.Error(), http.StatusInternalServerError)
		return
	}

	request.Reply(r, w, productRange, http.StatusOK)
}

func (a *Products) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := request.Param(r, "id")
	if id == "" {
		request.Reply(r, w, "missing 'id' parameter in path", http.StatusBadRequest)
	}

	product, err := a.store.Get(r.Context(), id)
	if err != nil {
		request.Reply(r, w, err.Error(), http.StatusInternalServerError)
		return
	}

	request.Reply(r, w, product, http.StatusOK)
}

func (a *Products) PutProduct(w http.ResponseWriter, r *http.Request) {
	id := request.Param(r, "id")

	product := types.Product{}
	if err := request.GetBody(w, r, &product); err != nil {
		request.Reply(r, w, err.Error(), http.StatusBadRequest)
		return
	}

	if product.Id != id {
		request.Reply(r, w, "id mismatch", http.StatusBadRequest)
		return
	}

	if err := a.store.Put(r.Context(), product); err != nil {
		request.Reply(r, w, err.Error(), http.StatusInternalServerError)
		return
	}

	request.Reply(r, w, &product, http.StatusOK)
}

func (a *Products) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := request.Param(r, "id")
	if err := a.store.Delete(r.Context(), id); err != nil {
		request.Reply(r, w, err.Error(), http.StatusOK)
		return
	}
	request.Reply(r, w, request.Result{Success: true}, http.StatusOK)
}
