.PHONY: sample test run

run:
	go run cmd/sample/main.go

sample:
	go build -o $@ cmd/$@/main.go

test:
	go test ./...tests
