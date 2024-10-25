install-protobuf:
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

install-grpc:
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

build-service:
	@cd ${PATTERN} && \
	protoc -I ./${SERVICE} \
		--go_out ./${SERVICE} \
		--go_opt paths=source_relative \
		--go-grpc_out ./${SERVICE} \
		--go-grpc_opt paths=source_relative \
		./${SERVICE}/*.proto
