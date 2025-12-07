all: test start

test:
	go test -v tests/*

start:
	go run cmd/main/main.go
