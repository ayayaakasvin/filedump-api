FROM golang:1.24.2 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o binary /app/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/config ./config
COPY --from=builder /app/binary .

RUN chmod +x ./binary

EXPOSE 8080

CMD [ "/app/binary" ]

# app itself requires a postgres database to run and cache as well as a redis instance for caching
# make sure to set the environment variables for database connection and redis cache