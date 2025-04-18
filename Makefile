# Makefile

.PHONY: image force-image build

bin := telee
src := $(wildcard *.go)

# Default target
${bin}: Makefile ${src}
	go build -v -o "${bin}"

# Docker targets
image:
	docker build -t ${USER}/telee .

force-image:
	docker build --no-cache -t ${USER}/telee .

.PHONY: goreleaser-build
goreleaser-build:
	goreleaser release --snapshot --clean

.PHONY: test
test:
	go test -v -race ./cmd/main.go
