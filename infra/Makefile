.PHONY: preview up destroy refresh stack-init test

STACK ?= dev

preview:
	pulumi preview --stack $(STACK)

up:
	pulumi up --stack $(STACK)

destroy:
	pulumi destroy --stack $(STACK)

refresh:
	pulumi refresh --stack $(STACK)

stack-init:
	pulumi stack init $(STACK)

test:
	go test ./...
