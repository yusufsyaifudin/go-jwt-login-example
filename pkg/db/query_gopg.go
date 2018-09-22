package db

import (
	"time"

	"fmt"

	"github.com/go-pg/pg"
	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	"github.com/golang-migrate/migrate/source/go_bindata"
	"github.com/rs/zerolog/log"
	"github.com/yusufsyaifudin/go-jwt-login-example/assets/migrations"
)

var logger = log.With().Str("pkg", "db").Logger()
var goPgConnection *pg.DB

// NewGoPgQuery will create new connection and returns 3 output,
// 1. connection to database, this should not be used other than to close the connection
// 2. implementation of db.Query interface where you can query with opened connection
// 3. error if any error occurred
func NewGoPgQuery(config *Config) (dbConn *pg.DB, query Query, err error) {
	dbOptions, err := pg.ParseURL(config.ConnectionString)
	if err != nil {
		return
	}

	dbOptions.PoolSize = 10
	dbOptions.IdleTimeout = time.Duration(5) * time.Second

	dbConn = pg.Connect(dbOptions)
	goPgConnection = dbConn

	if config.Debug {
		dbConn.OnQueryProcessed(func(event *pg.QueryProcessedEvent) {
			query, err := event.FormattedQuery()
			if err != nil {
				log.Printf("error when log query, %s", err.Error())
				return
			}

			elapsedTime := float64(time.Since(event.StartTime).Nanoseconds()) / float64(1000000)
			logger.Debug().
				Str("elapsedTime", fmt.Sprintf("%0.2f ms", elapsedTime)).
				Str("query", query).
				Msg("")
		})
	}

	// using implemented interface
	query = &QueryGoPg{
		config: config,
	}
	return
}

// QueryGoPg implements Query interface with github.com/go-pg/pg connection
type QueryGoPg struct {
	config *Config
}

// Raw will query to Postgres using raw sql and map the result into dst.
func (q *QueryGoPg) Raw(dst interface{}, sql string, args ...interface{}) (err error) {
	_, err = goPgConnection.Query(dst, sql, args...)
	return
}

// Exec will do query to Postgres without returning values
func (q *QueryGoPg) Exec(sql string, args ...interface{}) (err error) {
	_, err = goPgConnection.Exec(sql, args...)
	return
}

func (q *QueryGoPg) Migrate() error {
	s := bindata.Resource(
		migrations.AssetNames(),
		func(name string) ([]byte, error) {
			return migrations.Asset(name)
		},
	)

	d, err := bindata.WithInstance(s)
	if err != nil {
		logger.Error().Msgf("bindata instance fail %s", err.Error())
		return err
	}

	m, err := migrate.NewWithSourceInstance(
		"go-bindata",
		d,
		q.config.ConnectionString,
	)
	if err != nil {
		logger.Error().Msgf("fail do migration to host %s => %s", q.config.ConnectionString, err.Error())
		return err
	}

	// run your migrations and handle the errors above of course
	// if migration error, it will flag as dirty, and you must run it manually
	version, dirty, _ := m.Version()
	if dirty {
		err := m.Force(int(version))
		if err != nil {
			return err
		}

		return m.Up()
	}

	return m.Up()
}
