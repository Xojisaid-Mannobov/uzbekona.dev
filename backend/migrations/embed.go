// Package migrations SQL migratsiya fayllarini binary ichiga joylaydi.
package migrations

import "embed"

// FS — barcha *.sql fayllar (0001_name.up.sql / 0001_name.down.sql).
//
//go:embed *.sql
var FS embed.FS
