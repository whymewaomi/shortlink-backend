package core_config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Host             string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	RedisHost        string
	RedisPassword    string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		Host:             os.Getenv("HOST"),
		RedisHost:        os.Getenv("REDIS_HOST"),
		RedisPassword:    os.Getenv("REDIS_PASSWORD"),
		PostgresUser:     os.Getenv("POSTGRES_USER"),
		PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
		PostgresDB:       os.Getenv("POSTGRES_DB"),
	}
}

func (c *Config) Validate() error {
	cfg := map[string]string{
		"HOST":       c.Host,
		"REDIS_HOST": c.RedisHost,
	}

	for _, v := range cfg {
		if v == "" {
			panic("error config")
		}
	}

	return nil
}
