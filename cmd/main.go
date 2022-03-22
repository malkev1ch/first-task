package main

import (
	"context"
	"github.com/jackc/pgx/v4"
	server "github.com/malkev1ch/first-task"
	"github.com/malkev1ch/first-task/configs"
	"github.com/malkev1ch/first-task/internal/handler"
	"github.com/malkev1ch/first-task/internal/repository"
	"github.com/malkev1ch/first-task/internal/repository/postgres"
	"github.com/malkev1ch/first-task/internal/service"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := configs.Init("./configs")
	if err != nil {
		logrus.Fatal(err, "wrong configs variables")
	}

	db, err := newPostgresDB(cfg)
	if err != nil {
		logrus.Fatal(err, "err initializing DB")
	}

	repo := repository.NewRepositoryPostgres(db)
	services := service.NewService(repo)
	handlers := handler.NewHandler(services)
	srv := server.NewServer(cfg, handlers.InitRoutes())

	go func() {
		if err := srv.Run(); err != http.ErrServerClosed {
			logrus.Fatal(err, "error occurred while running http server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := srv.Shutdown(ctx); err != nil {
		logrus.Fatal(err,"failed to stop server")
	}

	if err := db.Close(context.Background()); err != nil {
		logrus.Fatal(err,"failed to stop connection db")
	}

}

func newPostgresDB(cfg *configs.Config) (*pgx.Conn, error) {
	return postgres.NewPostgresDB(postgres.Config{
		Host:     cfg.Postgres.Host,
		Port:     cfg.Postgres.Port,
		Username: cfg.Postgres.Username,
		Password: cfg.Postgres.Password,
		DBName:   cfg.Postgres.Dbname,
		SSLMode:  cfg.Postgres.Sslmode,
	})
}