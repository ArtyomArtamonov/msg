proto-compile: ./proto/messages.proto
	protoc --go_out=./pkg/api --go_opt=paths=source_relative --go-grpc_out=./pkg/api --go-grpc_opt=paths=source_relative ./proto/messages.proto
	protoc --go_out=./pkg/auth --go_opt=paths=source_relative --go-grpc_out=./pkg/auth --go-grpc_opt=paths=source_relative ./proto/auth.proto
