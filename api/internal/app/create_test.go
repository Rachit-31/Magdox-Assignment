package app

import (
	"encoding/json"
	"net/http"
	"testing"
)

// These two tests belong to ticket F1 and fail until the endpoint exists.
// They are a starting point, not the whole of the ticket.

func TestCreateProduct(t *testing.T) {
	forget(t, "KC-CUBE")
	c := signIn(t, "asha@kestrel.example")
	status, body := c.post("/api/products", map[string]string{"name": "  Kestrel Cube  ", "sku": "KC-CUBE", "deviceType": "camera"})
	if status != http.StatusCreated {
		t.Fatalf("status %d: %s", status, body)
	}
	var p Product
	if err := json.Unmarshal(body, &p); err != nil {
		t.Fatal(err)
	}
	if p.ID == "" || p.Name != "Kestrel Cube" || p.SKU != "KC-CUBE" || p.DeviceType != "camera" {
		t.Fatalf("created %+v", p)
	}
}

func TestViewerCannotCreate(t *testing.T) {
	forget(t, "KC-VIEW")
	in := map[string]string{"name": "Kestrel View", "sku": "KC-VIEW", "deviceType": "camera"}
	if status, _ := signIn(t, "meera@kestrel.example").post("/api/products", in); status != http.StatusForbidden {
		t.Fatalf("status %d, want 403", status)
	}
}

func TestCreateProductRejectsInvalidCSRF(t *testing.T) {
	c := signIn(t, "asha@kestrel.example")
	status, _ := c.do("POST", "/api/products", map[string]string{"name": "Kestrel Cube", "sku": "KC-CSRF", "deviceType": "camera"}, map[string]string{"X-CSRF-Token": "wrong"})
	if status != http.StatusForbidden {
		t.Fatalf("status %d, want 403", status)
	}
}

func TestCreateProductRejectsInvalidFields(t *testing.T) {
	c := signIn(t, "asha@kestrel.example")
	cases := []map[string]string{
		{"name": " ", "sku": "KC-EMPTY", "deviceType": "camera"},
		{"name": string(make([]byte, 0)), "sku": "kc-lower", "deviceType": "camera"},
		{"name": "Kestrel", "sku": "KC_BAD", "deviceType": "camera"},
		{"name": "Kestrel", "sku": "KC-TYPE", "deviceType": "phone"},
	}
	for _, in := range cases {
		if status, _ := c.post("/api/products", in); status != http.StatusBadRequest {
			t.Errorf("input %+v: status %d, want 400", in, status)
		}
	}
}

func TestCreateProductDuplicateSKUConflicts(t *testing.T) {
	const sku = "KC-DUPLICATE"
	forget(t, sku)
	c := signIn(t, "asha@kestrel.example")
	in := map[string]string{"name": "Kestrel Duplicate", "sku": sku, "deviceType": "camera"}
	if status, body := c.post("/api/products", in); status != http.StatusCreated {
		t.Fatalf("first create status %d: %s", status, body)
	}
	if status, _ := c.post("/api/products", in); status != http.StatusConflict {
		t.Fatalf("second create status %d, want 409", status)
	}
}

func TestCreateProductCanTargetAnotherOrganisation(t *testing.T) {
	const sku = "TW-NEW-1"
	forget(t, sku)
	c := signIn(t, "asha@kestrel.example")
	in := map[string]string{"name": "Tidewater Product", "sku": sku, "deviceType": "router", "orgId": "22222222-2222-4222-8222-222222222222"}
	if status, body := c.post("/api/products", in); status != http.StatusCreated {
		t.Fatalf("status %d: %s", status, body)
	}
	var orgID string
	if err := testDB.QueryRow(`SELECT org_id FROM products WHERE sku = $1`, sku).Scan(&orgID); err != nil {
		t.Fatal(err)
	}
	if orgID != "22222222-2222-4222-8222-222222222222" {
		t.Fatalf("product created in organisation %s", orgID)
	}
}

func TestAPIResponsesIncludeBuildHeader(t *testing.T) {
	status, _ := newClient(t).get("/api/products")
	if status != http.StatusUnauthorized {
		t.Fatalf("status %d, want 401", status)
	}
	res, err := http.Get(testServer.URL + "/api/unknown")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if got := res.Header.Get("X-Build"); got != "tidewater-7" {
		t.Fatalf("X-Build = %q, want tidewater-7", got)
	}
}
