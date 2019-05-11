package fruits

import (
	"database/sql"
	"strings"

	_ "github.com/lib/pq"
	"github.com/philip-bui/fruits-service/config"
	"github.com/rs/zerolog/log"
)

var (
	// DB is the exported PostgreSQL Client.
	DB *sql.DB
)

func init() {
	initDB()
}

func initDB() {
	if DB != nil {
		return
	}
	db, err := sql.Open("postgres", "host="+config.PostgresHost+
		" port="+config.PostgresPort+
		" user="+config.PostgresUser+
		" password="+config.PostgresPass+
		" dbname="+config.PostgresDB+
		" sslmode=disable")
	if err != nil {
		log.Fatal().Err(err).Msg("error connecting to postgres")
	}
	db.SetMaxIdleConns(0)
	DB = db
}

func PrepareStatement(query string) *sql.Stmt {
	initDB()
	stmt, err := DB.Prepare(query)
	if err != nil {
		log.Fatal().Err(err).Str("query", strings.ReplaceAll(strings.ReplaceAll(query, "\n", " "), "\t", "")).Msg("error preparing statement")
	}
	return stmt
}
