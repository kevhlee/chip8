SOURCES := main.go $(wildcard chip8/*.go)

bin/chip8: $(SOURCES)
	@ mkdir -p bin
	@ go build -o $@

clean:
	@ rm -rf bin
