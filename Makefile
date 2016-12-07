get-deps:
	go get ./...
	go get github.com/golang/mock/mockgen

s3mock.go:
	mockgen -source $(GOPATH)/src/github.com/aws/aws-sdk-go/service/s3/s3iface/interface.go -destination s3mock.go -package main

kmsmock.go:
	mockgen -source $(GOPATH)/src/github.com/aws/aws-sdk-go/service/kms/kmsiface/interface.go -destination kmsmock.go -package main

test: s3mock.go kmsmock.go
	go test -v ./...

install:
	go build && mv mayoiga $(GOPATH)/bin/mayoiga

.PHONY: get-deps test install
