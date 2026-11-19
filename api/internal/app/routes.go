package app

import (
	"net/http"
	"strings"
)

// Handler returns the API's routes.
func (a *App) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/session", a.handleSignIn)
	mux.HandleFunc("GET /api/products", a.withSession(a.handleProducts))
	mux.HandleFunc("POST /api/products", a.withSession(a.handleCreateProduct))
	mux.HandleFunc("GET /api/debug/rows", a.handleDebugRows)
	mux.HandleFunc("GET /api/findings", a.withSession(a.handleFindings))
	mux.HandleFunc("GET /api/export", a.withSession(a.handleExport))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api" || strings.HasPrefix(r.URL.Path, "/api/") {
			w.Header().Set("X-Build", "tidewater-7")
		}
		mux.ServeHTTP(w, r)
	})
}
