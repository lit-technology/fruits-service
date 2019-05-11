package main

import (
	"database/sql"
	"flag"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
	"github.com/philip-bui/fruits-service/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	Hash = "hash"
)

func main() {
	db, err := sql.Open("postgres", "host="+config.PostgresHost+
		" port="+config.PostgresPort+
		" user="+config.PostgresUser+
		" password="+config.PostgresPass+
		" dbname="+config.PostgresDWDB+
		" sslmode=disable")
	if err != nil {
		log.Fatal().Err(err).Msg("error connecting to Postgres")
	}
	var (
		p int
		t string
		q int
	)
	flag.IntVar(&p, "p", 2, "Number of Partitions")
	flag.StringVar(&t, "t", "answer", "Table name")
	flag.IntVar(&q, "q", 63, "Number of questions")
	flag.Parse()

	logError := func() *zerolog.Event {
		return log.Error().Int("partitions", p).Str("tableName", t).Int("questions", q)
	}
	if p > 2048 || p <= 0 {
		logError().Msg("unexpected partitions, expected 1 - 2048")
	} else if q > 64 || q <= 0 {
		logError().Msg("unexpected questions, expected 1 - 64")
	} else if len(t) == 0 {
		logError().Msg("invalid table name")
	}
	for i := 0; i < p; i++ {
		s := AnswerTable(i, q, t)
		//s := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %[1]s_%002[2]d PARTITION OF \"%[1]s\" FOR VALUES WITH (MODULUS %[3]d, REMAINDER %[2]d)", t, i, p)
		if _, err := db.Exec(s); err != nil {
			log.Fatal().Err(err).Msg("error creating partition")
		}
	}
}

func AnswerTable(i int, q int, t string) string {
	sb := strings.Builder{}
	for j := 1; j < q; j++ {
		sb.WriteString(fmt.Sprintf("N_%d SMALLINT,", j))
		sb.WriteString(fmt.Sprintf("NN_%d SMALLINT ARRAY,", j))
		sb.WriteString(fmt.Sprintf("S_%d VARCHAR(800),", j))
	}
	sb.WriteString(fmt.Sprintf("N_%d SMALLINT,", q))
	sb.WriteString(fmt.Sprintf("NN_%d SMALLINT ARRAY,", q))
	sb.WriteString(fmt.Sprintf("S_%d VARCHAR(800)", q))
	return fmt.Sprintf(`
				CREATE TABLE IF NOT EXISTS %s_%04d (
					survey_id BIGINT NOT NULL,
					user_id BIGINT,
					ip INET,
					user_agent VARCHAR(80),
					referrer VARCHAR(80),
					created TIMESTAMP NOT NULL DEFAULT NOW(),
					%s
				)
			`, t, i, sb.String())
}
