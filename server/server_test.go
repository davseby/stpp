package server

import (
	"database/sql"
	"fmt"
	"foodie/db"
	"net"
	"net/url"
	"os"
	"testing"

	"github.com/ory/dockertest/v3"
	dc "github.com/ory/dockertest/v3/docker"
)

var dbFn func(t *testing.T) *sql.DB

func TestMain(m *testing.M) {
	cleanup := setupDB()

	code := m.Run()

	if cleanup != nil {
		cleanup()
	}

	os.Exit(code)
}

// cleanUpTables must be called after each test function to clean up the
// database.
func cleanUpTables(t *testing.T, dbh *sql.DB) {
	t.Cleanup(func() {
		dbh.Exec("DELETE FROM plans")
		dbh.Exec("DELETE FROM recipes")
		dbh.Exec("DELETE FROM products")
		dbh.Exec("DELETE FROM users")
	})
}

func setupDB() (cleanup func()) {
	var err error

	dbFn = func(t *testing.T) *sql.DB {
		t.Fatal(err)
		return nil
	}

	var pool *dockertest.Pool

	pool, err = dockertest.NewPool("")
	if err != nil {
		return
	}

	var resource *dockertest.Resource

	resource, err = pool.Run("mariadb", "latest", []string{
		"MYSQL_ALLOW_EMPTY_PASSWORD=yes",
		"MYSQL_DATABASE=testing",
	})
	if err != nil {
		return
	}

	dsn := fmt.Sprintf(
		"root:@(%s)/testing?multiStatements=true",
		getHostPort(resource, "3306/tcp"),
	)

	var dbh *sql.DB

	if err = pool.Retry(func() error {
		dbh, err = db.Connect(dsn)
		return err
	}); err != nil {
		pool.Purge(resource)
		return
	}

	if err = db.Migrate(dbh); err != nil {
		pool.Purge(resource)
		return
	}

	dbh.Exec("DELETE FROM plans")
	dbh.Exec("DELETE FROM recipes")
	dbh.Exec("DELETE FROM products")
	dbh.Exec("DELETE FROM users")

	dbFn = func(t *testing.T) *sql.DB {
		return dbh
	}

	return func() {
		dbh.Close()
		pool.Purge(resource)
	}
}

func getHostPort(r *dockertest.Resource, portID string) string {
	dockerURL := os.Getenv("DOCKER_HOST")
	if dockerURL == "" {
		if r.Container == nil || r.Container.NetworkSettings == nil {
			return ""
		}

		m, ok := r.Container.NetworkSettings.Ports[dc.Port(portID)]
		if !ok || len(m) == 0 {
			return ""
		}

		ip := m[0].HostIP
		if ip == "0.0.0.0" || ip == "" {
			ip = "localhost"
		}

		return net.JoinHostPort(ip, m[0].HostPort)
	}

	u, err := url.Parse(dockerURL)
	if err != nil {
		panic(err)
	}

	return u.Hostname() + ":" + r.GetPort(portID)
}
