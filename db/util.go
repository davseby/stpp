package db

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
)

//go:embed sql
var migrations embed.FS

var (
	ErrNotFound  = errors.New("not found")
	ErrDuplicate = errors.New("duplicate")
)

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

func Migrate(dsn string) error {
	src, err := httpfs.New(http.FS(migrations), "sql")
	if err != nil {
		return err
	}
	if src == nil {
		return errors.New("couldn't load sql migrations file system")
	}

	m, err := migrate.NewWithSourceInstance("source", src, "mysql://"+dsn)
	if err != nil {
		return fmt.Errorf("creating migrations instance: %v", err)
	}

	err = m.Up()
	switch err {
	case nil:
		// Migration performed successfully.
	case migrate.ErrNoChange:
		// Schema is up to date.
	case os.ErrNotExist:
		// Schema is in unknown state, usually happens after application
		// roll-back when schema is newer than application expected.
	default:
		return err
	}

	errSource, errDB := m.Close()
	if errSource != nil {
		return errSource
	}
	if errDB != nil {
		return errDB
	}

	return nil
}
