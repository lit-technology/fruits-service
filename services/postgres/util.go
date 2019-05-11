package postgres

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
	"github.com/philip-bui/fruits-service/config"
	"github.com/rs/zerolog/log"
)

var (
	// DB is the exported PostgreSQL Client.
	DB         *sql.DB
	NullInt64  = sql.NullInt64{}
	NullBool   = sql.NullBool{}
	NullString = sql.NullString{}
)

func NewDB(n string) *sql.DB {
	if DB != nil {
		return DB
	}
	db, err := sql.Open("postgres", "host="+config.PostgresHost+
		" port="+config.PostgresPort+
		" user="+config.PostgresUser+
		" password="+config.PostgresPass+
		" dbname="+n+
		" sslmode=disable")
	if err != nil {
		log.Fatal().Err(err).Msg("error connecting to postgres")
	}
	db.SetMaxIdleConns(0)
	return db
}

func StringArrayNullable(a []string, sep string) interface{} {
	if a == nil || len(a) == 0 {
		return NullString
	}
	s := strings.Join(a, sep)
	return StringNullable(s)
}

func StringNullable(s string) interface{} {
	if s == "" {
		return NullString
	}
	return s
}

func Int64(i int64) *sql.NullInt64 {
	v := &sql.NullInt64{}
	v.Scan(i)
	return v
}

func Int64Nullable(i int64) interface{} {
	if i == 0 {
		return NullInt64
	}
	return i
}

func BoolNullable(b bool) interface{} {
	if !b {
		return NullBool
	}
	return b
}

func PrepareStatement(query string) *sql.Stmt {
	return nil
	if DB == nil {
		DB = NewDB(config.PostgresDB)
	}
	stmt, err := DB.Prepare(query)
	if err != nil {
		log.Fatal().Err(err).Str("query", strings.Replace(query, "\t", "", 100)).Msg("error preparing statement")
	}
	return stmt
}

func NewVarArgs(start int, end int) string {
	if end-start <= 0 {
		return ""
	}
	sb := strings.Builder{}
	for start < end-1 {
		sb.WriteString(fmt.Sprintf("$%d, ", start))
		start++
	}
	sb.WriteString(fmt.Sprintf("$%d", start))
	return sb.String()
}
