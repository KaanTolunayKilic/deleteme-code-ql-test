//go:generate go run github.com/prisma/prisma-client-go generate
package main

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"git.haw-hamburg.de/aci822/gabi/tweet-extractor/internal"
	"git.haw-hamburg.de/aci822/gabi/tweet-extractor/internal/config"
	"git.haw-hamburg.de/aci822/gabi/tweet-extractor/internal/utils"
)

func main() {
	logger := utils.NewLogger()
	config := config.ReadConfig(logger)
	app := internal.NewApp(logger, config)
	app.Initialize()
	_ = connect()
	app.Start()

	// Clean shutdown
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	<-sigc
	app.Shutdown()
	logger.Info("Exiting application after termination signal")
	os.Exit(0)
}

const (
	user     = "dbuser"
	password = "s3cretp4ssword"
)

func connect() *sql.DB {
	secretPassword := md5.Sum([]byte(password))
	connStr := fmt.Sprintf("postgres://%s:%s@localhost/pqgotest", user, secretPassword)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil
	}
	return db
}
