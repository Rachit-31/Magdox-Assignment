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
	status, body := c.post("/api/products", map[string]string{"name": "Kestrel Cube", "sku": "KC-CUBE", "deviceType": "camera"})
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
