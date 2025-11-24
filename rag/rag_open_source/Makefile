version ?= v1.1.0

docker-image-rag-arm64:
	docker build -f ./Dockerfile.arm64v8 --build-arg WANWU_ARCH=arm64 -t wanwulite/rag:$(version)-$(shell git rev-parse --short HEAD) .

docker-image-rag-amd64:
	docker build -f ./Dockerfile.amd64 --build-arg WANWU_ARCH=amd64 -t wanwulite/rag:$(version)-$(shell git rev-parse --short HEAD) .