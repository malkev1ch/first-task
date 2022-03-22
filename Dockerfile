FROM golang:1.17-alpine3.15 AS builder

COPY . /first-task/
WORKDIR /first-task/

RUN go mod download
RUN DB_PASSWORD=qwerty go build -o ./bin/app cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /first-task/bin/app .

EXPOSE 80

CMD ["./app"]