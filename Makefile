.PHONY: sample test run

run:
	go run cmd/sample/main.go

.SECONDEXPANSION:
sample process exact: cmd/$$@/main.go
	go build -o $@ cmd/$@/main.go

test:
	go test ./...tests
