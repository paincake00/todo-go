package app

import (
	"fmt"

	"github.com/paincake00/todo-go/internal/utils/env"
)

type Config struct {
	addr string
	db   DBConfig
}

type DBConfig struct {
	address string
	driver  string
}

func LoadConfig() Config {
	cfg := Config{
		addr: env.GetString("APP_ADDR", ":8080"),
		db: DBConfig{
			address: getPostgresUri(),
			driver:  env.GetString("DB_DRIVER", "postgres"),
		},
	}
	return cfg
}

func getPostgresUri() string {
	schema := env.GetString("DB_DRIVER", "postgres")
	user := env.GetString("POSTGRES_USER", "postgres")
	password := env.GetString("POSTGRES_PASSWORD", "postgres")
	host := env.GetString("POSTGRES_HOST", "localhost")
	port := env.GetString("POSTGRES_PORT", "5432")
	name := env.GetString("POSTGRES_DB", "todo-db")
	return fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable", schema, user, password, host, port, name)
}
