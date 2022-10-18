package db

import (
	"database/sql"
	"embed"
	"errors"
	"net/http"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	mmysql "github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
)

//go:embed sql
var migrations embed.FS

// ErrNotFound is returned whenever an object in the database is not found.
var ErrNotFound = errors.New("not found")

// Connect tries to establish a connection to the database by the provided
// dsn string. Once the connection is established, it returns an API to
// communicate with the database.
func Connect(dsn string) (*sql.DB, error) {
	cfg, err := mysql.ParseDSN(dsn)
	if err != nil {
		return nil, err
	}

	if cfg.Params == nil {
		cfg.Params = make(map[string]string)
	}

	cfg.Params["time_zone"] = `"+00:00"`
	cfg.ParseTime = true

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

// Migrate applies available migrations from sql folder to the
// given database handler.
func Migrate(dbh *sql.DB) error {
	src, err := httpfs.New(http.FS(migrations), "sql")
	if err != nil {
		return err
	}

	ins, err := mmysql.WithInstance(dbh, &mmysql.Config{})
	if err != nil {
		return err
	}

	migr, err := migrate.NewWithInstance("source", src, "mysql", ins)
	if err != nil {
		return err
	}

	err = migr.Up()
	switch err {
	case nil, migrate.ErrNoChange:
		// OK.
	case os.ErrNotExist:
		// Schema is in unknown state.
	default:
		return err
	}

	if err := src.Close(); err != nil {
		return err
	}

	return nil
}
