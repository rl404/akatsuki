package swagger

import (
	"github.com/go-chi/chi"
	httpSwagger "github.com/swaggo/http-swagger"
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
