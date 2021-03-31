pb:
	protoc --go_opt=paths=source_relative --go_out=. --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/userpb/*.proto

cleanpb:
	find . -name '*.pb.go' -delete