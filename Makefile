
GOPATH:=$(shell go env GOPATH)
.PHONY: proto
proto:
	protoc -I proto --go_out proto --go_opt paths=source_relative --go-grpc_out proto --go-grpc_opt paths=source_relative --grpc-gateway_out proto --grpc-gateway_opt paths=source_relative proto/pb/*.proto

.PHONY: build
build:
	go build -o micro *.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: vendor
vendor:
	go get ./...
	go mod vendor
	go mod verify

.PHONY: config
config:
	cp -rf ./config.example.yaml ./config.yaml
