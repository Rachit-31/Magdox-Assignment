package app

import (
	"context"
	"net/http"
)

type Finding struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Severity  string `json:"severity"`
	Component string `json:"component"` // empty when the component is not identified
	Product   string `json:"product"`
}

// findings returns the organisation's findings, the most severe first, then
// by title.
func (a *App) findings(ctx context.Context, s *Session) ([]Finding, error) {
	rows, err := a.org(s).Query(ctx,
		`SELECT f.id, f.title, f.severity, COALESCE(f.component, ''), p.name
		 FROM findings f JOIN products p ON p.id = f.product_id
		 WHERE f.org_id = $1
		 ORDER BY CASE f.severity WHEN 'critical' THEN 0 WHEN 'high' THEN 1 WHEN 'medium' THEN 2 ELSE 3 END, f.title, f.id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Finding{}
	for rows.Next() {
		var f Finding
		if err := rows.Scan(&f.ID, &f.Title, &f.Severity, &f.Component, &f.Product); err != nil {
			return nil, err
		}
		items = append(items, f)
	}
	return items, rows.Err()
}

func (a *App) handleFindings(w http.ResponseWriter, r *http.Request, s *Session) {
	items, err := a.findings(r.Context(), s)
	if err != nil {
		serverError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}
