package app

import (
	"fmt"
	"net/http"
)

// handleExport sends the organisation's findings as a file.
func (a *App) handleExport(w http.ResponseWriter, r *http.Request, s *Session) {
	format := r.URL.Query().Get("format")
	if format != "json" && format != "csv" {
		writeError(w, http.StatusBadRequest, "format must be json or csv")
		return
	}
	items, err := a.findings(r.Context(), s)
	if err != nil {
		serverError(w, err)
		return
	}
	w.Header().Set("Content-Disposition", `attachment; filename="findings.`+format+`"`)
	switch format {
	case "json":
		writeJSON(w, http.StatusOK, items)
	case "csv":
		w.Header().Set("Content-Type", "text/csv; charset=utf-8")
		fmt.Fprintln(w, "id,title,severity,component,product")
		for _, f := range items {
			fmt.Fprintf(w, "%s,%s,%s,%s,%s\n", f.ID, f.Title, f.Severity, f.Component, f.Product)
		}
	}
}
