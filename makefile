apidev:
	@air -c .air.api.toml

workerdev:
	@air -c .air.worker.toml

protobuf:
	@echo "Compiling protobuf..."
	@protoc --proto_path=protobuf --go_out=. --go-grpc_out=. protobuf/*.proto

.PHONY: apidev workerdev protobuf