package middleware

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/newrelic/go-agent/v3/newrelic"
)

// NewHTTP to create new newrelic http middleware.
func NewHTTP(nrApp *newrelic.Application) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get route pattern.
			path, ok := getRoutePattern(r)
			if !ok {
				path = r.RequestURI
			}

			// Start newrelic transaction.
			tx := nrApp.StartTransaction(r.Method + " " + path)
			defer tx.End()

			w = tx.SetWebResponse(w)
			tx.SetWebRequestHTTP(r)

			r = newrelic.RequestWithTransactionContext(r, tx)

			next.ServeHTTP(w, r)
		})
	}
}

func getRoutePattern(r *http.Request) (string, bool) {
	routePath := r.URL.Path
	if r.URL.RawPath != "" {
		routePath = r.URL.RawPath
	}

	rctx := chi.RouteContext(r.Context())
	tctx := chi.NewRouteContext()
	if rctx.Routes.Match(tctx, r.Method, routePath) {
		return tctx.RoutePattern(), true
	}

	return "", false
}
