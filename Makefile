
GOPATH:=$(shell go env GOPATH)
.PHONY: proto
proto:
	protoc --proto_path=. --go_out=plugins=grpc:. proto/*.proto
	
.PHONY: build
build:
	go build -o auth *.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: vendor
vendor:
	go get ./...
	go mod vendor
	go mod verify
