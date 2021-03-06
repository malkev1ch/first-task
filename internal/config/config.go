package config

type Config struct {
	CurrentDB           string `env:"CURRENT_DB" envDefault:"postgres"`
	PostgresURL         string `env:"POSTGRES_URL"`
	MongoURL            string `env:"MONGO_URL"`
	RedisURL            string `env:"REDIS_URL"`
	ImagePath           string `env:"IMAGE_PATH"`
	HTTPServer          string `env:"HTTP_SERVER_ADDRESS" envDefault:"localhost:8080"`
	JWTKey              string `env:"JWT_KEY"`
	CatsStreamName      string `env:"CATS_STREAM_NAME" envDefault:"cats"`
	CacheWorkersNum     int    `env:"CACHE_WORKERS_NUM" envDefault:"3"`
	CatsStreamGroupName string `env:"CATS_CONSUMERS_GROUP_NAME" envDefault:"consumers"`
	AuthMode            bool   `env:"AUTH_MODE" envDefault:"true"`
}
