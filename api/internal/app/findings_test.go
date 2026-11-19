package app

import (
	"encoding/json"
	"net/http"
	"testing"
)

func findingsOf(t *testing.T, c *client) []Finding {
	t.Helper()
	status, body := c.get("/api/findings")
	if status != http.StatusOK {
		t.Fatalf("status %d: %s", status, body)
	}
	var out struct {
		Items []Finding `json:"items"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		t.Fatal(err)
	}
	return out.Items
}

func TestFindingsAreListedMostSevereFirst(t *testing.T) {
	items := findingsOf(t, signIn(t, "asha@kestrel.example"))
	if len(items) != 12 {
		t.Fatalf("got %d findings, want 12", len(items))
	}
	if items[0].Severity != "critical" || items[len(items)-1].Severity != "low" {
		t.Fatalf("order: first %s, last %s", items[0].Severity, items[len(items)-1].Severity)
	}
	if n := len(findingsOf(t, signIn(t, "ravi@tidewater.example"))); n != 3 {
		t.Fatalf("Tidewater sees %d findings, want 3", n)
	}
}
