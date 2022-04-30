proto-c:
	protoc --go_out=./internal/message --go_opt=paths=source_relative --go-grpc_out=./internal/message --go-grpc_opt=paths=source_relative ./proto/message.proto
	protoc --go_out=./internal/auth --go_opt=paths=source_relative --go-grpc_out=./internal/auth --go-grpc_opt=paths=source_relative ./proto/auth.proto
