# One command for every check: `make check`. CI runs the same.
.PHONY: check fmt-check vet lint test imports build run preview

check: fmt-check vet lint imports test

fmt-check:
	@out="$$(gofmt -l cmd internal)"; if [ -n "$$out" ]; then echo "gofmt needed:"; echo "$$out"; exit 1; fi

vet:
	go vet ./...

lint:
	golangci-lint run ./...

imports:
	scripts/check-imports.sh

test:
	go test -count=1 ./...

build:
	CGO_ENABLED=0 go build -o bin/officeapp ./cmd/officeapp

# Local dev server on a fresh demo database in ./pb_data (delete it to start over).
run: build
	@[ -f pb_data/data.db ] || bin/officeapp seed --dir pb_data
	bin/officeapp serve --dir pb_data --http 127.0.0.1:8090

preview:
	scripts/preview.sh
