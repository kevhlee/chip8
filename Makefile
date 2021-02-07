##
## Build
##

.PHONY: build
build:
	go build -o ./bin/ch8 ./cmd/ch8

.PHONY: clean
clean:
	rm -rf bin

.PHONY: test
test:
	go test ./pkg/...
