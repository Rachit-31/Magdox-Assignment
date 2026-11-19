package app

import (
	"context"
	"net/http"
	"testing"

	"devshelf/internal/store"
)

func TestSignIn(t *testing.T) {
	c := signIn(t, "asha@kestrel.example")
	if c.csrf == "" {
		t.Fatal("sign in returned no CSRF token")
	}
	if status, _ := c.get("/api/products"); status != http.StatusOK {
		t.Fatalf("signed-in list: status %d", status)
	}
}

func TestUnknownUserCannotSignIn(t *testing.T) {
	status, _ := newClient(t).do("POST", "/api/session", map[string]string{"email": "nobody@nowhere.example"}, nil)
	if status != http.StatusUnauthorized {
		t.Fatalf("status %d, want 401", status)
	}
}

func TestListsNeedASession(t *testing.T) {
	for _, path := range []string{"/api/products", "/api/findings", "/api/export?format=json"} {
		if status, _ := newClient(t).get(path); status != http.StatusUnauthorized {
			t.Errorf("GET %s without a session: status %d, want 401", path, status)
		}
	}
}

func TestMigrateAndSeedCanRunAgain(t *testing.T) {
	ctx := context.Background()
	if err := store.Migrate(ctx, testDB); err != nil {
		t.Fatal(err)
	}
	if err := store.Seed(ctx, testDB); err != nil {
		t.Fatal(err)
	}
	var n int
	if err := testDB.QueryRowContext(ctx, `SELECT count(*) FROM orgs`).Scan(&n); err != nil {
		t.Fatal(err)
	}
	if n != 2 {
		t.Fatalf("%d organisations after seeding again, want 2", n)
	}
}
