run:
	GOOS=linux GOARCH=amd64 go run *.go

test:
	GOOS=linux GOARCH=amd64 go test ./...