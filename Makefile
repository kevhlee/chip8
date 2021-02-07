SOURCES := main.go $(wildcard src/*.go)

bin/chip8: $(SOURCES)
	@ mkdir -p bin
	@ go build -o $@

clean:
	@ rm -rf bin
