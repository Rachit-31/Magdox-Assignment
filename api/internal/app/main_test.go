package app

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"

	"devshelf/internal/store"
)

// The tests use their own database, emptied and filled again on every run.
var (
	testDB     *sql.DB
	testServer *httptest.Server
)

func TestMain(m *testing.M) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://devshelf@127.0.0.1:54329/devshelf_test?sslmode=disable"
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("the tests need the practice database (docker compose up -d): %v", err)
	}
	if _, err := db.ExecContext(ctx, `DROP SCHEMA public CASCADE; CREATE SCHEMA public`); err != nil {
		log.Fatal(err)
	}
	if err := store.Migrate(ctx, db); err != nil {
		log.Fatal(err)
	}
	if err := store.Seed(ctx, db); err != nil {
		log.Fatal(err)
	}
	testDB = db
	testServer = httptest.NewServer(New(db).Handler())
	code := m.Run()
	testServer.Close()
	db.Close()
	os.Exit(code)
}

// client is one browser: it keeps the session cookie and the CSRF token.
type client struct {
	t    *testing.T
	http *http.Client
	csrf string
}

func newClient(t *testing.T) *client {
	t.Helper()
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	return &client{t: t, http: &http.Client{Jar: jar}}
}

// signIn returns a client signed in as one of the made-up users.
func signIn(t *testing.T, email string) *client {
	t.Helper()
	c := newClient(t)
	status, body := c.do("POST", "/api/session", map[string]string{"email": email}, nil)
	if status != http.StatusOK {
		t.Fatalf("sign in as %s: status %d: %s", email, status, body)
	}
	var out struct {
		CSRFToken string `json:"csrfToken"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		t.Fatal(err)
	}
	c.csrf = out.CSRFToken
	return c
}

// do sends one request and returns the status and the body.
func (c *client) do(method, path string, body any, header map[string]string) (int, []byte) {
	c.t.Helper()
	var reader io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			c.t.Fatal(err)
		}
		reader = bytes.NewReader(raw)
	}
	req, err := http.NewRequest(method, testServer.URL+path, reader)
	if err != nil {
		c.t.Fatal(err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	for k, v := range header {
		req.Header.Set(k, v)
	}
	res, err := c.http.Do(req)
	if err != nil {
		c.t.Fatal(err)
	}
	defer res.Body.Close()
	raw, err := io.ReadAll(res.Body)
	if err != nil {
		c.t.Fatal(err)
	}
	return res.StatusCode, raw
}

func (c *client) get(path string) (int, []byte) { return c.do("GET", path, nil, nil) }

// post sends a change with the session's CSRF token.
func (c *client) post(path string, body any) (int, []byte) {
	return c.do("POST", path, body, map[string]string{"X-CSRF-Token": c.csrf})
}

type productPage struct {
	Items    []Product `json:"items"`
	Page     int       `json:"page"`
	PageSize int       `json:"pageSize"`
	HasNext  bool      `json:"hasNext"`
}

func products(t *testing.T, c *client, path string) productPage {
	t.Helper()
	status, body := c.get(path)
	if status != http.StatusOK {
		t.Fatalf("GET %s: status %d: %s", path, status, body)
	}
	var page productPage
	if err := json.Unmarshal(body, &page); err != nil {
		t.Fatal(err)
	}
	return page
}

// forget removes the products a test added when the test ends, so no test
// depends on which test ran before it.
func forget(t *testing.T, skus ...string) {
	t.Helper()
	t.Cleanup(func() {
		for _, sku := range skus {
			if _, err := testDB.Exec(`DELETE FROM products WHERE sku = $1`, sku); err != nil {
				t.Errorf("clean up %s: %v", sku, err)
			}
		}
	})
}
