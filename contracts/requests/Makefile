REPO_DIR=$(abspath ../../)

generate-go:
	mkdir -p langs/go && \
	protoc \
		-I . \
		-I .. \
		-I ../.. \
		--go-grpc_out=require_unimplemented_servers=false,paths=source_relative:langs/go \
		--go_out=paths=source_relative:langs/go \
		*.proto
	@echo "Client and swagger is generated"
