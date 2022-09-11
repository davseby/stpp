package main

import (
	"errors"
	"flag"
	"foodie/db"
	"foodie/server"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

func main() {
	var (
		dsn    string
		port   string
		secret string
	)

	flag.StringVar(&port, "port", "13307", "Server port")
	flag.StringVar(&dsn, "db", "root:db_password@tcp(127.0.0.1:13306)/db", "Database DSN")
	flag.StringVar(&secret, "secret", "sadghi21849adgjhlh904h3u4", "JWT secret")
	flag.Parse()

	dbh, err := db.Connect(dsn)
	if err != nil {
		logrus.WithError(err).
			Fatal("cannot connect to the database")
	}

	if err := db.Migrate(dsn); err != nil {
		logrus.WithError(err).
			Fatal("cannot apply migrations to the database")
	}

	srv := server.NewServer(dbh, port, []byte(secret))

	serverStop := make(chan struct{})
	go func() {
		logrus.Info("started web server")

		if err := srv.Start(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				logrus.WithError(err).Error("unexpected web server closure")
			}
		}

		close(serverStop)
	}()

	terminationCh := make(chan os.Signal, 1)
	signal.Notify(terminationCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-terminationCh:
	case <-serverStop:
	}

	if err := srv.Stop(); err != nil {
		logrus.WithError(err).Error("stopping web server")
	}

	<-serverStop

	logrus.Info("stopped web server")
}
