package user

import (
	"github.com/go-chi/chi"
)

// Router organization
func Router() chi.Router {
	r := chi.NewRouter()

	r.Get("/", list)
	r.Post("/", create)
	r.Route("/{permission_id}", func(r chi.Router) {
		r.Delete("/", delete)
	})

	return r
}
