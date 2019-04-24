
all: test build

test:
	go test -v ./objects/
	go test -v ./api/
#go test -v ./tests/common_test.go
#go test -v ./tests/subscribe_test.go
# go test -v ./tests/walletapi_test.go

build:
	go build ./...
