sample:
	go build -o $@ cmd/$@/main.go
test:
	go test ./...tests
