// Package store applies the migrations and loads the made-up practice data.
package store

import (
	"context"
	"database/sql"
	"fmt"
	"sort"

	"devshelf/migrations"
)

// Migrate applies every migration that has not been applied yet, in file name
// order, each in its own transaction. It is safe to run again.
func Migrate(ctx context.Context, db *sql.DB) error {
	if _, err := db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations (version text PRIMARY KEY)`); err != nil {
		return err
	}
	entries, err := migrations.FS.ReadDir(".")
	if err != nil {
		return err
	}
	var names []string
	for _, e := range entries {
		names = append(names, e.Name())
	}
	sort.Strings(names)
	for _, name := range names {
		var done bool
		if err := db.QueryRowContext(ctx, `SELECT EXISTS (SELECT 1 FROM schema_migrations WHERE version = $1)`, name).Scan(&done); err != nil {
			return err
		}
		if done {
			continue
		}
		body, err := migrations.FS.ReadFile(name)
		if err != nil {
			return err
		}
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, string(body)); err != nil {
			tx.Rollback()
			return fmt.Errorf("migration %s: %w", name, err)
		}
		if _, err := tx.ExecContext(ctx, `INSERT INTO schema_migrations (version) VALUES ($1)`, name); err != nil {
			tx.Rollback()
			return err
		}
		if err := tx.Commit(); err != nil {
			return err
		}
	}
	return nil
}

// The two made-up organisations.
const (
	KestrelID   = "11111111-1111-4111-8111-111111111111"
	TidewaterID = "22222222-2222-4222-8222-222222222222"
)

type product struct{ org, name, sku, deviceType string }

type finding struct{ org, sku, title, severity, component string }

var products = []product{
	{KestrelID, "Kestrel Bullet", "KC-BULLET", "camera"},
	{KestrelID, "Kestrel Bullet Pro", "KC-BULLET-PRO", "camera"},
	{KestrelID, "Kestrel Dome 2", "KC-DOME-2", "camera"},
	{KestrelID, "Kestrel Dome 4K", "KC-DOME-4K", "camera"},
	{KestrelID, "Kestrel Doorbell", "KC-DOORBELL", "camera"},
	{KestrelID, "Kestrel Fisheye", "KC-FISHEYE", "camera"},
	{KestrelID, "Kestrel Mini", "KC-MINI", "camera"},
	{KestrelID, "Kestrel NVR 16", "KC-NVR-16", "other"},
	{KestrelID, "Kestrel NVR 8", "KC-NVR-8", "other"},
	{KestrelID, "Kestrel PTZ 30", "KC-PTZ-30", "camera"},
	{KestrelID, "Kestrel Thermal", "KC-THERMAL", "camera"},
	{KestrelID, "Kestrel Turret", "KC-TURRET", "camera"},
	{TidewaterID, "Tidewater Edge 5", "TW-EDGE-5", "gateway"},
	{TidewaterID, "Tidewater Home Router", "TW-HOME-1", "router"},
	{TidewaterID, "Tidewater LTE Gateway", "TW-LTE-1", "gateway"},
	{TidewaterID, "Tidewater Mesh Node", "TW-MESH-1", "router"},
}

// An empty component means the engine could not say which program the
// finding belongs to; it is stored as NULL.
var findings = []finding{
	{KestrelID, "KC-DOME-2", "Telnet service enabled at boot", "critical", "telnetd-lite"},
	{KestrelID, "KC-DOME-2", "Default admin password, shared across devices", "critical", ""},
	{KestrelID, "KC-MINI", "Firmware image is not signed", "critical", "bootwarden"},
	{KestrelID, "KC-BULLET", `Web page echoes the "name" field without escaping`, "high", "fernhttpd"},
	{KestrelID, "KC-BULLET", "Update check uses HTTP & accepts any certificate", "high", "otafetchd"},
	{KestrelID, "KC-PTZ-30", "Length check uses < where <= is needed", "high", "linkcfgd"},
	{KestrelID, "KC-PTZ-30", "Debug shell reachable on the serial port", "medium", ""},
	{KestrelID, "KC-NVR-8", "Recordings folder is world readable", "medium", "clipstored"},
	{KestrelID, "KC-NVR-8", "Session cookie has no expiry", "medium", "fernhttpd"},
	{KestrelID, "KC-DOORBELL", "UPnP relay opens ports without asking", "low", "upnp-relay"},
	{KestrelID, "KC-FISHEYE", "Build date, host name and user left in the image", "low", ""},
	{KestrelID, "KC-THERMAL", "Old TLS versions still accepted", "low", "libcrypt-lite"},
	{TidewaterID, "TW-HOME-1", "Admin page reachable from the WAN side", "high", "fernhttpd"},
	{TidewaterID, "TW-MESH-1", "DNS relay answers for any name, from any source", "medium", "namerelayd"},
	{TidewaterID, "TW-LTE-1", "Config backup stores the Wi-Fi key in clear text", "critical", ""},
}

// Seed loads the practice data once. It does nothing when data is already there.
func Seed(ctx context.Context, db *sql.DB) error {
	var n int
	if err := db.QueryRowContext(ctx, `SELECT count(*) FROM orgs`).Scan(&n); err != nil {
		return err
	}
	if n > 0 {
		return nil
	}
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	steps := []struct {
		query string
		args  []any
	}{
		{`INSERT INTO orgs (id, name) VALUES ($1, 'Kestrel Cameras'), ($2, 'Tidewater Routers')`, []any{KestrelID, TidewaterID}},
		{`INSERT INTO users (org_id, email, name, role) VALUES
			($1, 'asha@kestrel.example', 'Asha Verma', 'analyst'),
			($1, 'meera@kestrel.example', 'Meera Nair', 'viewer'),
			($2, 'ravi@tidewater.example', 'Ravi Menon', 'analyst')`, []any{KestrelID, TidewaterID}},
	}
	for _, s := range steps {
		if _, err := tx.ExecContext(ctx, s.query, s.args...); err != nil {
			return err
		}
	}
	for _, p := range products {
		if _, err := tx.ExecContext(ctx, `INSERT INTO products (org_id, name, sku, device_type) VALUES ($1, $2, $3, $4)`,
			p.org, p.name, p.sku, p.deviceType); err != nil {
			return err
		}
	}
	for _, f := range findings {
		if _, err := tx.ExecContext(ctx, `INSERT INTO findings (org_id, product_id, title, severity, component)
			SELECT $1, id, $3, $4, NULLIF($5, '') FROM products WHERE org_id = $1 AND sku = $2`,
			f.org, f.sku, f.title, f.severity, f.component); err != nil {
			return err
		}
	}
	return tx.Commit()
}
