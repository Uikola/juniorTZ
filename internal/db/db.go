package db

import (
	"database/sql"
	"github.com/Uikola/juniorTZ/internal/config"
	"github.com/rs/zerolog"
	"time"
)

func InitDB(cfg *config.Config, log zerolog.Logger) *sql.DB {
	var db *sql.DB
	var err error
	for i := 0; i < 3; i++ {
		db, err = sql.Open(cfg.DriverName, cfg.ConnString)
		if err != nil {
			log.Error().Err(err).Msg("failed to connect to the database. Retrying in 5 seconds")
			time.Sleep(5 * time.Second)
			continue
		}

		err = db.Ping()
		if err != nil {
			log.Error().Err(err).Msg("failed to ping the database. Retrying in 5 seconds")
			time.Sleep(5 * time.Second)
			continue
		}

		break
	}

	return db
}
