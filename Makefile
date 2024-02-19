build:
	go build ./cmd/main.go

build-linux:
	GOOS=linux GOARCH=amd64 go build ./cmd/main.go
