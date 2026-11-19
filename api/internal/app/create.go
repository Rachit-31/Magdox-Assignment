package app

import (
	"crypto/subtle"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/jackc/pgx/v5/pgconn"
)

func validSKU(s string) bool {
	if len(s) < 3 || len(s) > 20 {
		return false
	}
	for _, r := range s {
		if (r < 'A' || r > 'Z') && (r < '0' || r > '9') && r != '-' {
			return false
		}
	}
	return true
}

func (a *App) handleCreateProduct(w http.ResponseWriter, r *http.Request, s *Session) {
	if subtle.ConstantTimeCompare([]byte(r.Header.Get("X-CSRF-Token")), []byte(s.CSRF)) != 1 {
		writeError(w, http.StatusForbidden, "invalid CSRF token")
		return
	}
	if s.Role != "analyst" && s.Role != "admin" {
		writeError(w, http.StatusForbidden, "insufficient permissions")
		return
	}
	var in struct {
		Name       string `json:"name"`
		SKU        string `json:"sku"`
		DeviceType string `json:"deviceType"`
		OrgID      string `json:"orgId"`
	}
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<16))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	in.Name = strings.TrimSpace(in.Name)
	if n := utf8.RuneCountInString(in.Name); n < 1 || n > 80 {
		writeError(w, http.StatusBadRequest, "name must be 1 to 80 characters")
		return
	}
	if !validSKU(in.SKU) {
		writeError(w, http.StatusBadRequest, "sku must be 3 to 20 uppercase letters, digits or hyphens")
		return
	}
	switch in.DeviceType {
	case "camera", "router", "gateway", "other":
	default:
		writeError(w, http.StatusBadRequest, "invalid deviceType")
		return
	}
	orgID := s.OrgID
	if in.OrgID != "" && in.OrgID != s.OrgID {
		var exists bool
		if err := a.DB.QueryRowContext(r.Context(), `SELECT EXISTS (SELECT 1 FROM orgs WHERE id::text = $1)`, in.OrgID).Scan(&exists); err != nil {
			serverError(w, err)
			return
		}
		if !exists {
			writeError(w, http.StatusBadRequest, "unknown organisation")
			return
		}
		orgID = in.OrgID
	}
	var p Product
	err := a.DB.QueryRowContext(r.Context(), `INSERT INTO products (org_id, name, sku, device_type)
		VALUES ($1, $2, $3, $4) RETURNING id, name, sku, device_type`, orgID, in.Name, in.SKU, in.DeviceType).
		Scan(&p.ID, &p.Name, &p.SKU, &p.DeviceType)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			writeError(w, http.StatusConflict, "sku already exists for this organisation")
			return
		}
		serverError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, p)
}
