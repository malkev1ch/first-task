package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
	"time"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(cfg Config) (*pgx.Conn, error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	db, err := pgx.Connect(ctx, fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		logrus.Fatal(err, "error creating connect to DB")
	}

	return db, nil
}
