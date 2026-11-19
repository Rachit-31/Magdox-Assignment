package app

import (
	"context"
	"testing"

	"devshelf/internal/store"
)

func TestProductStatusMigrationIsIdempotentAndDefaultsToActive(t *testing.T) {
	if err := store.Migrate(context.Background(), testDB); err != nil {
		t.Fatalf("rerun migrations: %v", err)
	}
	var nonActive int
	if err := testDB.QueryRow(`SELECT count(*) FROM products WHERE status <> 'active'`).Scan(&nonActive); err != nil {
		t.Fatal(err)
	}
	if nonActive != 0 {
		t.Fatalf("found %d existing products without active status", nonActive)
	}
	const sku = "F3-OLD-INSERT"
	forget(t, sku)
	if _, err := testDB.Exec(`INSERT INTO products (org_id, name, sku, device_type) VALUES ($1, $2, $3, $4)`, store.KestrelID, "Legacy Insert", sku, "other"); err != nil {
		t.Fatalf("insert without status: %v", err)
	}
	var status string
	if err := testDB.QueryRow(`SELECT status FROM products WHERE sku = $1`, sku).Scan(&status); err != nil {
		t.Fatal(err)
	}
	if status != "active" {
		t.Fatalf("legacy insert status = %q, want active", status)
	}
	if _, err := testDB.Exec(`UPDATE products SET status = 'paused' WHERE sku = $1`, sku); err == nil {
		t.Fatal("invalid product status was accepted")
	}
	if _, err := testDB.Exec(`UPDATE products SET status = NULL WHERE sku = $1`, sku); err == nil {
		t.Fatal("empty product status was accepted")
	}
}
