IMG ?= tyk-broadcast:test
CLUSTER_NAME ?= kind

.PHONY: cross-build-image
cross-build-image: ## Build docker image
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o manager.linux *.go
	docker build -f Dockerfile . -t ${IMG}
	go run hack/load.go -image ${IMG} -cluster=${CLUSTER_NAME}