##
## Build
##

.PHONY: build
build:
	go build -o ./bin/ch8 ./main.go

.PHONY: clean
clean:
	rm -rf bin

.PHONY: test
test:
	go test ./pkg/...

##
## Installation
##

.PHONY: install
install: build
	install -d /usr/local/bin
	install -m755 bin/ch8 /usr/local/bin

.PHONY: uninstall
uninstall:
	rm -f /usr/local/bin/ch8
