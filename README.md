# Msg
[![ci](https://github.com/ArtyomArtamonov/msg/actions/workflows/ci.yml/badge.svg)](https://github.com/ArtyomArtamonov/msg/actions/workflows/ci.yml)

Msg is gRPC-based backend for messaging written in golang.

## Goals

Future goals:

- Implement message service which will allow users to send messages to chat rooms
- Unit tests and github workflow

Already achieved:

- Authentication and authorization with JWT and refresh tokens
- PostgreSQL as persistance storage
- gRPC and protobuf service creation with help of Makefile and well-prepared code organization

## Compiling proto messages

Run

```
make proto-c
```

For that to work you need to install

```
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```

## Env variables

Create .env file and put it in a project root

.env file has to contain

```env
HOST="192.168.0.2:1234"

JWT_SECRET="s3cr3t"
JWT_DURATION_MIN=15
REFRESH_DURATION_DAYS=90


POSTGRES_DB="dbname"
POSTGRES_USER="dbuser"
POSTGRES_PASSWORD="dbuserpassword"

PGADMIN_DEFAULT_EMAIL="email"
PGADMIN_DEFAULT_PASSWORD="password"
PGADMIN_CONFIG_SERVER_MODE="False"
```

## Run

```console
$ docker-compose up -d
```
