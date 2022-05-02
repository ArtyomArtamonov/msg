proto-c:
	protoc --go_out=./internal/server --go_opt=paths=source_relative --go-grpc_out=./internal/server --go-grpc_opt=paths=source_relative ./proto/*.proto

build: cmd/main.go
	go build -o bin/program cmd/main.go
