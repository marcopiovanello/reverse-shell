windows-server:
	GOOS=windows GOARCH=amd64 go build -o reverse.exe -ldflags="-w -s" cmd/server-e2e/main.go