package app

import (
	"net/http"
	"strconv"
	"strings"
)

const pageSize = 5

type Product struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	SKU        string `json:"sku"`
	DeviceType string `json:"deviceType"`
}

// handleProducts lists the organisation's products by name, one page at a
// time. q keeps only the products whose name contains it, in any letter case.
func (a *App) handleProducts(w http.ResponseWriter, r *http.Request, s *Session) {
	q := strings.TrimSpace(r.URL.Query().Get("q"))
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
		 WHERE org_id = $1 AND name ILIKE '%` + q + `%'
		 ORDER BY name, id LIMIT $2 OFFSET $3`
	rows, err := a.org(s).Query(r.Context(), query, pageSize+1, (page-1)*pageSize)
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
	writeJSON(w, http.StatusOK, map[string]any{"items": items, "page": page, "pageSize": pageSize, "hasNext": hasNext})
}
