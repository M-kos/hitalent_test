package db

import (
	"context"
	"net"
	"net/url"

	"github.com/M-kos/hitalent_test/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const dBScheme = "postgres"

func OpenGormConnection(ctx context.Context, conf *config.Config) (*gorm.DB, error) {
	return gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn(conf),
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
}

func dsn(config *config.Config) string {
	query := url.Values{}
	query.Set("sslmode", "disable")

	value := url.URL{
		Scheme:   dBScheme,
		User:     url.UserPassword(config.Postgres.User, config.Postgres.Password),
		Host:     net.JoinHostPort(config.Postgres.Host, config.Postgres.Port),
		Path:     config.Postgres.Name,
		RawQuery: query.Encode(),
	}

	return value.String()
}
