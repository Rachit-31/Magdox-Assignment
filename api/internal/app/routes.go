package app

import "net/http"

// Handler returns the API's routes.
func (a *App) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/session", a.handleSignIn)
	mux.HandleFunc("GET /api/products", a.withSession(a.handleProducts))
	mux.HandleFunc("GET /api/findings", a.withSession(a.handleFindings))
	mux.HandleFunc("GET /api/export", a.withSession(a.handleExport))
	return mux
}
