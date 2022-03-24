run:
	POSTGRES_URL=postgres://postgres:qwerty@localhost:5433/postgres?sslmode=disable IMAGE_PATH=./Data/ HTTP_SERVER_ADDRESS=127.0.0.1:8080 go run main.go

build:
	POSTGRES_URL=postgres://postgres:qwerty@localhost:5432/postgres IMAGE_PATH=Data/CatImage/ go build -o ./bin cmd/main.go

image:
	docker build -t first-task-local-image:v1 .

container:
	docker run --name first-task-local -p 8080:8080 first-task-image:v1

postgres:
	docker run --name=first-task-local-db -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=qwerty -e POSTGRES_DB=postgres -p 5433:5432 -d postgres:14