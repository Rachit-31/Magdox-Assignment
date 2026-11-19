// Package app is the Devshelf HTTP API: sign-in, products, findings and exports.
package app

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

// App holds what every handler needs.
type App struct {
	DB *sql.DB

	mu       sync.Mutex
	sessions map[string]*Session
}

func New(db *sql.DB) *App {
	return &App{DB: db, sessions: map[string]*Session{}}
}

// orgDB runs statements for one organisation: its id is always the first
// parameter, $1.
type orgDB struct {
	db    *sql.DB
	orgID string
}

func (a *App) org(s *Session) orgDB { return orgDB{db: a.DB, orgID: s.OrgID} }

func (o orgDB) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return o.db.QueryContext(ctx, query, append([]any{o.orgID}, args...)...)
}

func (o orgDB) QueryRow(ctx context.Context, query string, args ...any) *sql.Row {
	return o.db.QueryRowContext(ctx, query, append([]any{o.orgID}, args...)...)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("write response: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

// serverError logs the cause and tells the caller nothing about it.
func serverError(w http.ResponseWriter, err error) {
	log.Printf("server error: %v", err)
	writeError(w, http.StatusInternalServerError, "something went wrong")
}
