.PHONY: build build-sdk build-all run dev test clean lint tidy web-dev web-build landing-rules

BINARY_NAME=booltools-seo-crawler
SDK_BINARY=seo-crawler

build:
	go build -o bin/$(BINARY_NAME) ./cmd/server

build-sdk:
	go build -o bin/$(SDK_BINARY) ./cmd/seo-crawler

build-all: build build-sdk

run: build
	./bin/$(BINARY_NAME)

dev:
	go run ./cmd/server

test:
	go test ./... -v -count=1

lint:
	go vet ./...

tidy:
	go mod tidy

web-dev:
	cd web && npm run dev

web-build:
	cd web && npm run build

landing-rules:
	cd landing && node generate-rule-pages.js

clean:
	go clean
	-rm -rf bin 2>/dev/null || rd /s /q bin 2>nul || true
	-rm -f data.db 2>/dev/null || del data.db 2>nul || true
