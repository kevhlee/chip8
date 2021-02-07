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

##
## Installation
##

DESTDIR :=
prefix  := /usr/local
bindir  := ${prefix}/bin

.PHONY: install
install: build
	install -d ${DESTDIR}${bindir}
	install -m755 bin/ch8 ${DESTDIR}${bindir}/

.PHONY: uninstall
uninstall:
	rm -f ${DESTDIR}${bindir}/ch8
