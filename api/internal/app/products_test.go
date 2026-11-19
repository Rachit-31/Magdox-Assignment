package app

import (
	"fmt"
	"net/http"
	"net/url"
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

func TestProductsPaginationReturnsEachProductOnce(t *testing.T) {
	const org = "Kestrel Cameras"
	var orgID string
	if err := testDB.QueryRow(`SELECT id FROM orgs WHERE name = $1`, org).Scan(&orgID); err != nil {
		t.Fatal(err)
	}
	skus := []string{"F2-PAGE-1", "F2-PAGE-2", "F2-PAGE-3", "F2-PAGE-4", "F2-PAGE-5", "F2-PAGE-6", "F2-PAGE-7"}
	forget(t, skus...)
	for i, sku := range skus {
		if _, err := testDB.Exec(`INSERT INTO products (org_id, name, sku, device_type) VALUES ($1, $2, $3, 'other')`, orgID, fmt.Sprintf("F2 Pagination %02d", i), sku); err != nil {
			t.Fatal(err)
		}
	}
	var expected int
	if err := testDB.QueryRow(`SELECT count(*) FROM products WHERE org_id = $1`, orgID).Scan(&expected); err != nil {
		t.Fatal(err)
	}
	c := signIn(t, "asha@kestrel.example")
	seen := map[string]bool{}
	for pageNumber := 1; ; pageNumber++ {
		page := products(t, c, fmt.Sprintf("/api/products?page=%d", pageNumber))
		if len(page.Items) > pageSize {
			t.Fatalf("page %d has %d items, want at most %d", pageNumber, len(page.Items), pageSize)
		}
		if page.HasNext != (pageNumber*pageSize < expected) {
			t.Fatalf("page %d hasNext=%v, want %v", pageNumber, page.HasNext, pageNumber*pageSize < expected)
		}
		for _, item := range page.Items {
			if seen[item.ID] {
				t.Fatalf("product %q appeared on more than one page", item.ID)
			}
			seen[item.ID] = true
		}
		if !page.HasNext {
			break
		}
	}
	if len(seen) != expected {
		t.Fatalf("saw %d products across pages, want %d", len(seen), expected)
	}
}

func TestSearchTreatsInputAsLiteralText(t *testing.T) {
	const sku = "F2-SEARCH-QUOTE"
	forget(t, sku)
	var orgID string
	if err := testDB.QueryRow(`SELECT org_id FROM users WHERE email = $1`, "asha@kestrel.example").Scan(&orgID); err != nil {
		t.Fatal(err)
	}
	if _, err := testDB.Exec(`INSERT INTO products (org_id, name, sku, device_type) VALUES ($1, $2, $3, 'other')`, orgID, "O'Brien %_ Camera", sku); err != nil {
		t.Fatal(err)
	}
	c := signIn(t, "asha@kestrel.example")
	for _, query := range []string{"O'Brien", "%_ Camera"} {
		page := products(t, c, "/api/products?q="+url.QueryEscape(query))
		found := false
		for _, item := range page.Items {
			if item.SKU == sku {
				found = true
			}
		}
		if !found {
			t.Fatalf("search %q did not find the matching product", query)
		}
	}
}

func TestDebugRowsDoesNotRequireSession(t *testing.T) {
	status, body := newClient(t).get("/api/debug/rows")
	if status != http.StatusOK || len(body) == 0 {
		t.Fatalf("GET /api/debug/rows: status %d, body %s", status, body)
	}
}
