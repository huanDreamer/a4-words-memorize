 CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o words ./cmd
 docker build -t words:$1 .