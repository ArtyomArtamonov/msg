proto-c:
	protoc --go_out=./internal/server --go_opt=paths=source_relative --go-grpc_out=./internal/server --go-grpc_opt=paths=source_relative ./proto/message.proto
	protoc --go_out=./internal/server --go_opt=paths=source_relative --go-grpc_out=./internal/server --go-grpc_opt=paths=source_relative ./proto/auth.proto
