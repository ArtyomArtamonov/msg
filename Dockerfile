########### PRODUCTION ###########
FROM golang:1.18-alpine as prod

WORKDIR /app

COPY . ./

RUN go mod download

RUN go build -v -o /bin/program ./cmd/main.go

RUN cp .env /.env

EXPOSE 50051

CMD ["/bin/program"]

########### DEBUG ###########
FROM golang:1.18 as debug

WORKDIR /app

COPY . ./

RUN go mod download
RUN go get github.com/go-delve/delve/cmd/dlv
RUN go install github.com/go-delve/delve/cmd/dlv

RUN go build -race -v -o /bin/program ./cmd/main.go

RUN cp .env /.env

EXPOSE 50051 2345

COPY ./dlv.sh /
RUN chmod +x /dlv.sh
ENTRYPOINT ["/dlv.sh"]
