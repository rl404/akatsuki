package ping

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rl404/akatsuki/internal/utils"
)

// Ping contains basic routes.
type Ping struct{}

// New to create new ping and other base routes.
func New() *Ping {
	return &Ping{}
}

// Register to register common routes.
func (p Ping) Register(r chi.Router) {
	r.Get("/", p.handleRoot)
	r.Get("/ping", p.handlePing)
	r.Get("/favicon.ico", p.handleFavIcon)
	r.Get("/rl404", p.handlerl404)
	r.NotFound(http.HandlerFunc(p.handleNotFound))
	r.MethodNotAllowed(http.HandlerFunc(p.handleMethodNotAllowed))
	r.Handle("/metrics", promhttp.Handler())
}

func (p Ping) handleRoot(w http.ResponseWriter, _ *http.Request) {
	utils.ResponseWithJSON(w, http.StatusOK, "ok", nil)
}

func (p Ping) handlePing(w http.ResponseWriter, _ *http.Request) {
	utils.ResponseWithJSON(w, http.StatusOK, "pong", nil)
}

func (p Ping) handleNotFound(w http.ResponseWriter, _ *http.Request) {
	utils.ResponseWithJSON(w, http.StatusNotFound, nil, nil)
}

func (p Ping) handleMethodNotAllowed(w http.ResponseWriter, _ *http.Request) {
	utils.ResponseWithJSON(w, http.StatusMethodNotAllowed, nil, nil)
}

func (p Ping) handleFavIcon(w http.ResponseWriter, _ *http.Request) {
	utils.ResponseWithJSON(w, http.StatusOK, "ok", nil)
}

func (p Ping) handlerl404(w http.ResponseWriter, _ *http.Request) {
	utils.ResponseWithJSON(w, http.StatusOK, "rl404 was here", nil)
}
