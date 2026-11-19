package app

import (
	"net/http"
	"testing"
)

func TestSearchFiltersByNameInAnyCase(t *testing.T) {
	page := products(t, signIn(t, "asha@kestrel.example"), "/api/products?q=nvr")
	if len(page.Items) != 2 {
		t.Fatalf("got %d products for q=nvr, want 2", len(page.Items))
	}
}

func TestEachOrganisationSeesItsOwnProducts(t *testing.T) {
	page := products(t, signIn(t, "ravi@tidewater.example"), "/api/products")
	if len(page.Items) != 4 || page.HasNext {
		t.Fatalf("Tidewater sees %d products (hasNext %v), want its 4", len(page.Items), page.HasNext)
	}
}

func TestBadPageIsRefused(t *testing.T) {
	c := signIn(t, "asha@kestrel.example")
	for _, p := range []string{"0", "-1", "abc"} {
		if status, _ := c.get("/api/products?page=" + p); status != http.StatusBadRequest {
			t.Errorf("page=%s: status %d, want 400", p, status)
		}
	}
}
