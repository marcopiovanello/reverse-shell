windows-server:
	GOOS=windows GOARCH=amd64 go build -o reverse.exe -ldflags="-w -s" cmd/server-e2e/main.go

windows-stager:
	GOOS=windows GOARCH=amd64 go build -o ciao.exe -ldflags="-w -s" cmd/stager/main.go

linux-client:
	GOOS=linux GOARCH=amd64 go build -o oculis -ldflags="-w -s" cmd/cli-e2e/main.go