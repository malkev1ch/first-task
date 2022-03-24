FROM golang:1.17-alpine3.15 AS builder

COPY . /github.com/malkev1ch/first-task/
WORKDIR /github.com/malkev1ch/first-task/

RUN go mod download
RUN go build -o ./bin/app main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /github.com/malkev1ch/first-task/bin/app .

EXPOSE 8080
CMD ["mkdir Data"]
CMD ["./app"]