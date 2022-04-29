.PHONY: test all

scripts:=sample process exact analyze


all: $(scripts)

run:
	go run cmd/sample/main.go

.SECONDEXPANSION:
$(scripts): cmd/$$@/main.go
	go build -o $@ cmd/$@/main.go

test:
	go test ./...tests
