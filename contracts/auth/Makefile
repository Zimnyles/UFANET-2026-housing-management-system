REPO_DIR=$(abspath ../../)

generate-go:
	mkdir -p langs/go && \
	docker run --platform linux/amd64 --rm \
		-v $(REPO_DIR):$(REPO_DIR) -v $(CURDIR):$(CURDIR) -w $(CURDIR) \
		artifactory.tutu.ru:443/docker-tutu/tutu-protoc:latest \
		-I . \
		-I .. \
		-I ../.. \
		--go-grpc_out=require_unimplemented_servers=false,paths=source_relative:langs/go \
		--go_out=paths=source_relative:langs/go \
		--grpc-gateway_out=allow_patch_feature=false,paths=source_relative:langs/go \
		--swagger_out=. \
		*.proto
	@echo "Client and swagger is generated"