run:
	DB_PASSWORD=qwerty go run cmd/main.go

image:
	docker build -t first-task-image:v1 .

container:
	docker run --name first-task -p 80:80 first-task-image:v1

postgres:
	docker run --name=first-task -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=qwerty -e POSTGRES_DB=postgres -p 5434:5432 -d postgres:14