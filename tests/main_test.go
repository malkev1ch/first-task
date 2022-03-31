package tests

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	log "github.com/sirupsen/logrus"
)

var golangPort string

func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// Create postgres container
	postgres, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "14",
		Env: []string{
			"POSTGRES_PASSWORD=qwerty",
			"POSTGRES_USER=postgres",
			"POSTGRES_DB=postgres",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start postgres: %s", err)
	}

	postgresHostAndPort := postgres.GetHostPort("5432/tcp")

	// Start migration
	flywayPath := fmt.Sprintf("-url=jdbc:postgresql://%s/postgres", postgresHostAndPort)
	flywayConfFile := "-configFiles=/home/andreimalkevich/first-task/sql/dockertest/flyway.conf"
	cmd := exec.Command("flyway", flywayPath, flywayConfFile, "migrate")
	if err := cmd.Run(); err != nil {
		log.Fatalf("Command finished with error: %v", err)
	}

	postgres.Expire(300) // Tell docker to hard kill the container in 300 seconds
	log.Info("Created postgres container successfully")

	postgresIPInNetworkBridge := postgres.GetIPInNetwork(&dockertest.Network{
		Network: &docker.Network{
			Name: "bridge",
		},
	})
	postgresURL := fmt.Sprintf("postgres://postgres:qwerty@%s:5432/postgres?sslmode=disable", postgresIPInNetworkBridge)

	// Create redis container
	redis, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "redis",
		Tag:        "7.0-rc-alpine",
		Env:        nil,
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start redis: %s", err)
	}

	redis.Expire(300) // Tell docker to hard kill the container in 300 seconds
	log.Info("Created redis container successfully")

	redisIPInNetworkBridge := redis.GetIPInNetwork(&dockertest.Network{
		Network: &docker.Network{
			Name: "bridge",
		},
	})
	redisURL := fmt.Sprintf("redis://:@%s:6379/1", redisIPInNetworkBridge)

	// pulls an image, creates a container based on it and runs it
	golang, err := pool.BuildAndRunWithOptions("/home/andreimalkevich/first-task/Dockerfile",
		&dockertest.RunOptions{
			Name: "test-handler",
			PortBindings: map[docker.Port][]docker.PortBinding{
				"8080/tcp": {docker.PortBinding{
					HostIP:   "0.0.0.0",
					HostPort: "8081",
				}},
			},
			Env: []string{
				fmt.Sprintf("POSTGRES_URL=%s", postgresURL),
				fmt.Sprintf("REDIS_URL=%s", redisURL),
				"MONGO_URL=mongodb://admin:qwerty@mongodb:27017/?maxPoolSize=20&w=majority",
				"IMAGE_PATH=/Data/",
				"HTTP_SERVER_ADDRESS=:8080",
				"CURRENT_DB=postgres",
				"JWT_KEY=secret_key_for_jwt",
				"CATS_STREAM_NAME=cats",
				"CATS_CONSUMERS_GROUP_NAME=consumers",
				"CACHE_WORKERS_NUM=3",
				"AUTH_MODE=false",
			},
		}, func(config *docker.HostConfig) {
			// set AutoRemove to true so that stopped container goes away by itself
			config.AutoRemove = true
			config.RestartPolicy = docker.RestartPolicy{Name: "no"}
		},
	)
	if err != nil {
		log.Fatalf("Could not start golang: %s", err)
	}

	golangPort = golang.GetPort("8080/tcp")
	golang.Expire(300) // Tell docker to hard kill the container in 300 seconds
	log.Infof("Created golang container successfully, port - %s", golangPort)
	time.Sleep(time.Minute * 2)
	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(postgres); err != nil {
		log.Fatalf("Could not purge postgres: %s", err)
	}

	if err := pool.Purge(redis); err != nil {
		log.Fatalf("Could not purge redis: %s", err)
	}

	if err := pool.Purge(golang); err != nil {
		log.Fatalf("Could not purge golang: %s", err)
	}
	os.Exit(code)
}
