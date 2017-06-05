TEST?= $(shell go list ./... | grep -v vendor)
DEPS = $(shell go list -f '{{range .TestImports}}{{.}} {{end}}' ./...)
BIN  = history

all: build

build: deps
	mkdir -p bin
	go build -o bin/$(BIN)

install: build
	install -m 755 ./bin/$(BIN) ~/bin/$(BIN)
	install -m 644 ./misc/zsh/completions/_history ~/.zsh/Completion

deps:
	go get -d -v ./...
	echo $(DEPS) | xargs -n1 go get -d

test: deps
	go test $(TEST) $(TESTARGS) -timeout=3s -parallel=4
	go vet $(TEST)
	go test $(TEST) -race

.PHONY: all build deps test
