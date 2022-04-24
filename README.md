# Msg

Msg is gRPC-based backend for messaging written in golang.

## Goals

Future goals:

- Authentication and authorization with JWT and refresh tokens
- PostgreSQL as persistance storage
- Implement message service which will allow users to send messages to chat rooms
- Unit tests and github workflow

Already achieved:

- JWT access token workflow
- gRPC and protobuf service creation with help of Makefile and well-prepared code organization
- UserStore intreface for further persistent storage implementation

## Env variables

Create .env file and put it in a project root

.env file has to contain

```env
HOST="192.168.0.1:1234"

JWT_SECRET="s3cr3t"
JWT_DURATION_MIN=15
```


