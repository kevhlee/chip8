##
## Build
##

.PHONY: build
build:
	go build -o ./bin/chip8 ./main.go

.PHONY: clean
clean:
	rm -rf bin

##
## Installation
##

.PHONY: install
install: build
	install -d /usr/local/bin
	install -m755 ./bin/chip8 /usr/local/bin

.PHONY: uninstall
uninstall:
	rm -f /usr/local/bin/chip8
