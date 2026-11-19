// Command devshelf runs the Devshelf API.
//
//	devshelf          start the API (the same as "devshelf serve")
//	devshelf migrate  apply the database migrations
//	devshelf seed     load the made-up practice data
package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"devshelf/internal/app"
	"devshelf/internal/store"
)

func env(name, fallback string) string {
	if v := os.Getenv(name); v != "" {
		return v
	}
	return fallback
}

func main() {
	db, err := sql.Open("pgx", env("DATABASE_URL", "postgres://devshelf@127.0.0.1:54329/devshelf?sslmode=disable"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	ctx := context.Background()
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("cannot reach the database (is it running? try: docker compose up -d): %v", err)
	}
	command := "serve"
	if len(os.Args) > 1 {
		command = os.Args[1]
	}
	switch command {
	case "migrate":
		if err := store.Migrate(ctx, db); err != nil {
			log.Fatal(err)
		}
		log.Print("migrations applied")
	case "seed":
		if err := store.Seed(ctx, db); err != nil {
			log.Fatal(err)
		}
		log.Print("practice data loaded")
	case "serve":
		addr := env("ADDR", "127.0.0.1:8080")
		server := &http.Server{Addr: addr, Handler: app.New(db).Handler(), ReadHeaderTimeout: 5 * time.Second}
		log.Printf("Devshelf API on http://%s", addr)
		log.Fatal(server.ListenAndServe())
	default:
		log.Fatalf("unknown command %q: use serve, migrate or seed", command)
	}
}
