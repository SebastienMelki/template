.PHONY: lint breaking deps

lint:
	cd .. && buf lint

breaking:
	cd .. && buf breaking --against '.git#branch=main'

deps:
	cd .. && buf dep update
