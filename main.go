//	Package classification of Cats storage API
//
//	Documentation for Cats storage API
//
//	Schemes: http
//	Host: localhost:8080
//	BasePath: /
//	Version: 1.0.0
//
//	Consumes:
//	 - application/json
//
//	Produces:
//	 - application/json
//	swagger:meta
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"github.com/malkev1ch/first-task/internal/config"
	"github.com/malkev1ch/first-task/internal/handler"
	"github.com/malkev1ch/first-task/internal/repository"
	"github.com/malkev1ch/first-task/internal/repository/mongodb"
	"github.com/malkev1ch/first-task/internal/repository/postgres"
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

	repo, err := CreateDBConnection(&cfg)
	if err != nil {
		logrus.Fatal(err, "err initializing DB")
	}

	services := service.NewService(repo)
	handlers := handler.NewHandler(services, &cfg)

	router := echo.New()
	router.Logger.SetLevel(log.DEBUG)
	router.Use(middleware.Logger())
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.POST, echo.DELETE},
	}))

	cat := router.Group("/cats")

	{
		cat.GET("/:uuid", handlers.GetCat)
		cat.POST("/", handlers.CreateCat)
		cat.PUT("/:uuid", handlers.UpdateCat)
		cat.DELETE("/:uuid", handlers.DeleteCat)
		cat.POST("/:uuid/image", handlers.UploadCatImage)
		cat.GET("/:uuid/image", handlers.GetCatImage)
	}
	//TODO JWT
	//auth := router.Group("/auth")
	//{
	//	auth.POST("/sign-up", handlers.signUp)
	//	auth.POST("/sign-in", handlers.signIn)
	//	auth.POST("/logout", handlers.signIn)
	//}

	router.Logger.Fatal(router.Start(cfg.HTTPServer))

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

func CreateDBConnection(cfg *config.Config) (repository.Repository, error) {
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
		return mongodb.RepositoryMongo{DB: client}, err
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
		return postgres.RepositoryPostgres{DB: conn}, err
	}

	logrus.WithFields(logrus.Fields{
		"status": "database connection failed.",
		"err":    "invalid config",
	}).Info("repository info")

	return nil, nil
}
