proto-c:
	protoc --go_out=./internal/server --go_opt=paths=source_relative --go-grpc_out=./internal/server --go-grpc_opt=paths=source_relative ./proto/*.proto

build: cmd/api_service/main.go cmd/message_service/main.go
	go build -o bin/api_service cmd/api_service/main.go
	go build -o bin/message_service cmd/message_service/main.go
