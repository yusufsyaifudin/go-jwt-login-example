package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/golang-migrate/migrate"
	"github.com/namsral/flag"
	"github.com/rs/zerolog/log"
	"github.com/yusufsyaifudin/go-jwt-login-example/pkg/auth"
	"github.com/yusufsyaifudin/go-jwt-login-example/pkg/db"
	"github.com/yusufsyaifudin/go-jwt-login-example/server"
)

var serverSecretKey = flag.String("secret-key", "ndjsHJUTUI8uok", "Server secret key")
var listenAddress = flag.String("listen-address", "localhost:8000", "Address to bind")
var dbUrl = flag.String("db-url", "postgres://postgres:postgres@localhost:5432/go-users?sslmode=disable", "Connection string to postgres")
var dbDebug = flag.Bool("db-debug", true, "Whether to show sql debug or not")
var logger = log.With().Str("pkg", "main").Logger()

// @title Authentication System
// @version 3.0
// @description This is a documentation for Authentication System
// @termsOfService http://example.com

// @contact.name API Support
// @contact.url http://www.example.com
// @contact.email contact.us@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8000
// @BasePath /api/v1
func main() {
	flag.Parse()

	dbConfig := &db.Config{
		ConnectionString: *dbUrl,
		Debug:            *dbDebug,
	}

	dbConnection, query, err := db.NewGoPgQuery(dbConfig)
	defer dbConnection.Close()
	if err != nil {
		logger.Error().Err(err).Msg("database connection fail")
		return
	}

	if err := query.Migrate(); err != nil && err != migrate.ErrNoChange {
		logger.Error().Err(err).Msg("migration fail")
	}

	srv := &server.Config{
		ListenAddress:   *listenAddress,
		ServerSecretKey: *serverSecretKey,
		DB:              query,
		Auth:            auth.NewJwtAuth(),
	}

	var apiErrChan = make(chan error, 1)
	go func() {
		logger.Info().Msgf("running api at %s", *listenAddress)
		apiErrChan <- srv.Run()
	}()

	// to gracefully shutdown the server
	var signalChan = make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	select {
	case <-signalChan:
		logger.Info().Msg("got an interrupt, exiting...")
		srv.Shutdown()
	case err := <-apiErrChan:
		if err != nil {
			logger.Error().Err(err).Msg("error while running api, exiting...")
		}
	}

}
