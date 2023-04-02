package context

import (
	"github.com/jmoiron/sqlx"
	// sqlite
	_ "github.com/mattn/go-sqlite3"
)

// NewDBConnection returns database connection for bellboy
func NewDBConnection(dbPath string) *sqlx.DB {
	return sqlx.MustConnect("sqlite3", dbPath)
}
