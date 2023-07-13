package swagger

import (
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// Swagger to handle swagger API docs.
type Swagger struct{}

// New to create new swagger route.
func New() *Swagger {
	return &Swagger{}
}

// Register to register swagger route.
func (sw *Swagger) Register(r chi.Router) {
	r.Get("/swagger/*", httpSwagger.WrapHandler)
}
