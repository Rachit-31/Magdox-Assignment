package app

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

const pageSize = 5

type Product struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	SKU        string `json:"sku"`
	DeviceType string `json:"deviceType"`
}

func (a *App) handleDebugRows(w http.ResponseWriter, r *http.Request) {
	rows, err := a.DB.QueryContext(r.Context(), `SELECT id, org_id, name, sku, device_type, created_at FROM products ORDER BY id`)
	if err != nil {
		serverError(w, err)
		return
	}
	defer rows.Close()
	type productRow struct {
		ID         string    `json:"id"`
		OrgID      string    `json:"orgId"`
		Name       string    `json:"name"`
		SKU        string    `json:"sku"`
		DeviceType string    `json:"deviceType"`
		CreatedAt  time.Time `json:"createdAt"`
	}
	items := []productRow{}
	for rows.Next() {
		var p productRow
		if err := rows.Scan(&p.ID, &p.OrgID, &p.Name, &p.SKU, &p.DeviceType, &p.CreatedAt); err != nil {
			serverError(w, err)
			return
		}
		items = append(items, p)
	}
	if err := rows.Err(); err != nil {
		serverError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

// handleProducts lists the organisation's products by name, one page at a
// time. q keeps only the products whose name contains it, in any letter case.
func (a *App) handleProducts(w http.ResponseWriter, r *http.Request, s *Session) {
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	q = strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace(q)
	page := 1
	if p := r.URL.Query().Get("page"); p != "" {
		n, err := strconv.Atoi(p)
		if err != nil || n < 1 {
			writeError(w, http.StatusBadRequest, "page must be a whole number from 1")
			return
		}
		page = n
	}
	query := `SELECT id, name, sku, device_type FROM products
		 WHERE org_id = $1 AND name ILIKE $2 ESCAPE '\'
		 ORDER BY name, id LIMIT $3 OFFSET $4`
	rows, err := a.org(s).Query(r.Context(), query, "%"+q+"%", pageSize+1, (page-1)*pageSize)
	if err != nil {
		serverError(w, err)
		return
	}
	defer rows.Close()
	items := []Product{}
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.SKU, &p.DeviceType); err != nil {
			serverError(w, err)
			return
		}
		items = append(items, p)
	}
	if err := rows.Err(); err != nil {
		serverError(w, err)
		return
	}
	hasNext := len(items) > pageSize
	if hasNext {
		items = items[:pageSize]
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items, "page": page, "pageSize": pageSize, "hasNext": hasNext})
}
