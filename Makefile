.PHONY: generate build test lint clean

generate:
	buf generate
	$(MAKE) -C backend generate-sqlc

build:
	$(MAKE) -C backend build

test:
	$(MAKE) -C backend test

lint:
	$(MAKE) -C contracts lint
	$(MAKE) -C backend lint

clean:
	$(MAKE) -C backend clean
