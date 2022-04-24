proto-c:
	protoc --go_out=./pkg/message --go_opt=paths=source_relative --go-grpc_out=./pkg/message --go-grpc_opt=paths=source_relative ./proto/message.proto
	protoc --go_out=./pkg/auth --go_opt=paths=source_relative --go-grpc_out=./pkg/auth --go-grpc_opt=paths=source_relative ./proto/auth.proto
