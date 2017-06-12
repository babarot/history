TEST?= $(shell go list ./... | grep -v vendor)
DEPS = $(shell go list -f '{{range .TestImports}}{{.}} {{end}}' ./...)
BIN  = history

all: build

build: deps
	mkdir -p bin
	go build -o bin/$(BIN)

install: build
	go install
	if echo $$SHELL | grep "zsh" &>/dev/null; then \
		install -m 644 ./misc/zsh/completions/_history $(shell zsh -c 'echo $$fpath[1]'); \
		fi

deps:
	go get -d -v ./...
	echo $(DEPS) | xargs -n1 go get -d

test: deps
	go test $(TEST) $(TESTARGS) -timeout=3s -parallel=4
	go vet $(TEST)
	go test $(TEST) -race

.PHONY: all build deps test
