package configs

type Config struct {
	CurrentDB   string `env:"CURRENTDB" envDefault:"postgres"`
	PostgresURL string `env:"POSTGRES_URL"`
	MongoURL    string `env:"MONGO_URL"`
	ImagePath   string `env:"IMAGE_PATH"`
	HTTPServer  string `env:"HTTP_SERVER_ADDRESS" envDefault:"localhost:8080"`
}
