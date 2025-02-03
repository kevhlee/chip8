NAME    := chip8
SOURCES := main.go $(wildcard chip8/*.go)

bin/$(NAME): $(SOURCES)
	@ mkdir -p bin
	@ go build -o $@

.PHONY: clean
clean:
	@ rm -rf bin

.PHONY: test
test:
	@ go test ./...
