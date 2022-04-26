FROM golang:1.18-alpine

WORKDIR /app

COPY . ./

RUN go mod download

RUN go build -v -o /bin/program ./cmd/main.go

RUN cp .env /.env

EXPOSE 50051

CMD ["/bin/program"]
