# Makefile

.PHONY: build
build:
	go build -trimpath -o ./tmp/telee ./cmd/main.go

.PHONY: test
test:
	go test -v -race ./cmd/main.go

.PHONY: release
release:
	git bump
	goreleaser release --rm-dist
