package dw

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/philip-bui/fruits-service/config"
	"github.com/rs/zerolog/log"
)

var (
	// DB is the exported PostgreSQL Client.
	DB *sql.DB
)

func init() {
	db, err := sql.Open("postgres", "host="+config.PostgresHost+
		" port="+config.PostgresPort+
		" user="+config.PostgresUser+
		" password="+config.PostgresPass+
		" dbname="+config.PostgresDWDB+
		" sslmode=disable")
	if err != nil {
		log.Fatal().Err(err).Msg("error connecting to postgres dw")
	}
	db.SetMaxIdleConns(0)
	DB = db
}
