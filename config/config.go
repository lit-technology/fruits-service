package config

import (
	"flag"

	"github.com/philip-bui/fruits-service/pkg/env"
)

const (
	DeltaTimestamp = 1555614095
)

var (
	// Port of gRPC Service.
	Port string
	// PostgresHost is host of PostgreSQL.
	PostgresHost string
	// PostgresPort is port of PostgreSQL.
	PostgresPort string
	// PostgresUser is user to connect to PostgreSQL.
	PostgresUser string
	// PostgresPass is password to connect to PostgreSQL.
	PostgresPass string
	// PostgresDB is default PostgreSQL Database to connect to.
	PostgresDB string
	// PostgreDWDB is data warehouse PostgreSQL Database to connect to.
	PostgresDWDB string

	JWT = []byte("")
)

func init() {
	flag.StringVar(&Port, "port", "8000", "PORT")
	flag.StringVar(&PostgresHost, "postgreshost", "127.0.0.1", "POSTGRES_HOST")
	flag.StringVar(&PostgresPort, "postgresport", "5432", "POSTGRES_PORT")
	flag.StringVar(&PostgresUser, "postgresuser", "postgres", "POSTGRES_USER")
	flag.StringVar(&PostgresPass, "postgrespass", "", "POSTGRES_PASS")
	flag.StringVar(&PostgresDB, "postgresdb", "fruits", "POSTGRES_DB")
	flag.StringVar(&PostgresDWDB, "postgresdwdb", "fruits_dw", "POSTGRES_DB")
	flag.Parse()

	env.LoadEnv(&Port, "PORT")
	env.LoadEnv(&PostgresHost, "POSTGRES_HOST")
	env.LoadEnv(&PostgresPort, "POSTGRES_PORT")
	env.LoadEnv(&PostgresUser, "POSTGRES_USER")
	env.LoadEnv(&PostgresPass, "POSTGRES_PASS")
	env.LoadEnv(&PostgresDB, "POSTGRES_DB")
	env.LoadEnv(&PostgresDWDB, "POSTGRES_DW_DB")
}
