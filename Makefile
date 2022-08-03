.PHONY: test install

install:
	go install ./...

update:
	go mod tidy

test:
	go test ./...tests
