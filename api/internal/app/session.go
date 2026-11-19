package app

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
)

const sessionCookie = "devshelf_session"

// Session is one signed-in person. It lives in memory: restarting the API
// signs everyone out, which is fine for a practice app.
type Session struct {
	Name    string
	Email   string
	Role    string // viewer, analyst or admin
	OrgID   string
	OrgName string
	CSRF    string
}

func randomToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// handleSignIn signs in one of the made-up users by email. There are no
// passwords in this practice app.
func (a *App) handleSignIn(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<16)).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	s := &Session{Email: in.Email}
	err := a.DB.QueryRowContext(r.Context(),
		`SELECT u.name, u.role, o.id, o.name FROM users u JOIN orgs o ON o.id = u.org_id WHERE u.email = $1`,
		in.Email).Scan(&s.Name, &s.Role, &s.OrgID, &s.OrgName)
	if errors.Is(err, sql.ErrNoRows) {
		writeError(w, http.StatusUnauthorized, "unknown user")
		return
	}
	if err != nil {
		serverError(w, err)
		return
	}
	token, err := randomToken()
	if err != nil {
		serverError(w, err)
		return
	}
	if s.CSRF, err = randomToken(); err != nil {
		serverError(w, err)
		return
	}
	a.mu.Lock()
	a.sessions[token] = s
	a.mu.Unlock()
	http.SetCookie(w, &http.Cookie{Name: sessionCookie, Value: token, Path: "/", HttpOnly: true, SameSite: http.SameSiteLaxMode})
	writeJSON(w, http.StatusOK, map[string]any{
		"user":      map[string]string{"name": s.Name, "email": s.Email, "role": s.Role, "org": s.OrgName},
		"csrfToken": s.CSRF,
	})
}

type sessionHandler func(http.ResponseWriter, *http.Request, *Session)

// withSession lets only signed-in callers through.
func (a *App) withSession(next sessionHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(sessionCookie)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "sign in first")
			return
		}
		a.mu.Lock()
		s := a.sessions[c.Value]
		a.mu.Unlock()
		if s == nil {
			writeError(w, http.StatusUnauthorized, "sign in first")
			return
		}
		next(w, r, s)
	}
}
