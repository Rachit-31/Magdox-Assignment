package app

import (
	"encoding/json"
	"net/http"
	"testing"
)

func exportJSON(t *testing.T, c *client) []Finding {
	t.Helper()
	status, body := c.get("/api/export?format=json")
	if status != http.StatusOK {
		t.Fatalf("json: status %d: %s", status, body)
	}
	var items []Finding
	if err := json.Unmarshal(body, &items); err != nil {
		t.Fatal(err)
	}
	return items
}

func TestExportJSON(t *testing.T) {
	if n := len(exportJSON(t, signIn(t, "asha@kestrel.example"))); n != 12 {
		t.Fatalf("got %d findings, want 12", n)
	}
}

func TestUnknownExportFormatIsRefused(t *testing.T) {
	if status, _ := signIn(t, "asha@kestrel.example").get("/api/export?format=pdf"); status != http.StatusBadRequest {
		t.Fatalf("status %d, want 400", status)
	}
}
