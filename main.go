// Cats storage API
//
// Documentation for Cats storage API
//
// Schemes: http
// Host: localhost:8080
// BasePath: /
// Version: 1.0.0
//
// Consumes:
//  - application/json
//
// Produces:
//  - application/json
//
// SecurityDefinitions:
//  AdminAuth:
//   type: apiKey
//   name: Authorization
//   in: header
//
// swagger:meta
package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/malkev1ch/first-task/internal/rediscache"

	"github.com/caarlos0/env/v6"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/malkev1ch/first-task/internal/config"
	"github.com/malkev1ch/first-task/internal/handler"
	"github.com/malkev1ch/first-task/internal/repository"
	"github.com/malkev1ch/first-task/internal/service"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	cfg := config.Config{}
	if err := env.Parse(&cfg); err != nil {
		logrus.Fatal(err, "wrong config variables")
	}
	logrus.Infof("Parsed config - %+v\n", cfg)
	repo, err := CreateDBConnection(&cfg)
	if err != nil {
		logrus.Fatal(err, "err initializing DB")
	}

	redisClient := redisConnection(cfg)
	defer func() {
		err := redisClient.Close()
		if err != nil {
			logrus.Errorf("error while closing redis connection - %e", err)
		}
	}()

	cache := rediscache.NewStreamCache(&cfg, redisClient)
	services := service.NewService(repo, cache)
	validator := handler.NewValidator()
	handlers := handler.NewHandler(services, &cfg, validator)
	router := handler.InitRouter(handlers, &cfg)

	// router.Logger.Fatal(router.Start(cfg.HTTPServer))

	router.Logger.Fatal(router.StartServer(&http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := router.Shutdown(ctx); err != nil {
		logrus.Fatal(err, "failed to stop server")
	}
}

func CreateDBConnection(cfg *config.Config) (*repository.Repository, error) {
	switch cfg.CurrentDB {
	case "mongo":
		client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.MongoURL))
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"status": "connection to mongodb database failed.",
				"err":    err,
			}).Fatal("mongodb repository info")
		} else {
			logrus.WithFields(logrus.Fields{
				"status": "successfully connected to mongodb database.",
			}).Info("mongodb repository info.")
		}
		return repository.NewRepositoryMongo(client), err
	case "postgres":
		conn, err := pgxpool.Connect(context.Background(), cfg.PostgresURL)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"status": "connection to postgres database failed.",
				"err":    err,
			}).Fatal("postgres repository info.")
		} else {
			logrus.WithFields(logrus.Fields{
				"status": "successfully connected to postgres database.",
			}).Info("postgres repository info.")
		}
		return repository.NewRepositoryPostgres(conn), err
	}

	logrus.WithFields(logrus.Fields{
		"status": "database connection failed.",
		"err":    "invalid config",
	}).Info("repository info")

	return nil, nil
}

func redisConnection(cfg config.Config) *redis.Client {
	opt, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"status":  "error while parsing connection URL for redis",
			"err":     err,
			"options": opt,
		}).Fatal("redis repository info.")
		return nil
	}
	redisClient := redis.NewClient(opt)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		logrus.WithFields(logrus.Fields{
			"status": "error while connection to redis",
			"err":    err,
		}).Fatal("redis repository info.")
		return nil
	}

	logrus.WithFields(logrus.Fields{
		"status": "successfully connected to redis",
	}).Info("redis repository info.")

	return redisClient
}
