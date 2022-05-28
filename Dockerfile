########### API PRODUCTION ###########
FROM golang:1.18-alpine as api_prod

WORKDIR /app

COPY . ./

RUN go mod download

RUN go build -v -o /bin/program ./cmd/api_service/main.go

RUN cp .env /.env

EXPOSE 50051

CMD ["/bin/program"]

########### API DEBUG ###########
FROM golang:1.18 as api_debug

WORKDIR /app

COPY . ./

RUN go mod download
RUN go get github.com/go-delve/delve/cmd/dlv
RUN go install github.com/go-delve/delve/cmd/dlv

RUN go build -v -o /bin/program ./cmd/api_service/main.go

RUN cp .env /.env

EXPOSE 50051 2345

COPY ./dlv.sh /
RUN chmod +x /dlv.sh
ENTRYPOINT ["/dlv.sh", "/app/cmd/api_service", "2345"]

########### MESSAGE PRODUCTION ###########
FROM golang:1.18-alpine as message_prod

WORKDIR /app

COPY . ./

RUN go mod download

RUN go build -v -o /bin/program ./cmd/message_service/main.go

RUN cp .env /.env

EXPOSE 50052

CMD ["/bin/program"]

########### MESSAGE DEBUG ###########
FROM golang:1.18 as message_debug

WORKDIR /app

COPY . ./

RUN go mod download
RUN go get github.com/go-delve/delve/cmd/dlv
RUN go install github.com/go-delve/delve/cmd/dlv

RUN go build -race -v -o /bin/program ./cmd/message_service/main.go

RUN cp .env /.env

EXPOSE 50052 2346

COPY ./dlv.sh /
RUN chmod +x /dlv.sh
ENTRYPOINT ["/dlv.sh", "/app/cmd/message_service" , "2346"]
