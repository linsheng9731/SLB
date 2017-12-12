
# do not specify a full path for go since travis will fail
GO = GOGC=off go
GOVENDOR = $(shell which govendor)
VENDORFMT = $(shell which vendorfmt)

all: build test

build: gofmt
	$(GO) build -i
	$(GO) test -i ./...

test: gofmt
	$(GO) test -v -test.timeout 15s `go list ./... | grep -v '/vendor/'`

checkdeps:
	[ -x "$(GOVENDOR)" ] || $(GO) get -u github.com/kardianos/govendor
	govendor list +e | grep '^ e ' && { echo "Found missing packages. Please run 'govendor add +e'"; exit 1; } || : echo

vendorfmt:
	[ -x "$(VENDORFMT)" ] || $(GO) get -u github.com/magiconair/vendorfmt/cmd/vendorfmt
	vendorfmt

gofmt:
	gofmt -w `find . -type f -name '*.go' | grep -v vendor`

install:
	$(GO) install

clean:
	$(GO) clean
	rm -rf pkg

.PHONY: build buildpkg clean docker gofmt homebrew install linux pkg release test vendorfmt
