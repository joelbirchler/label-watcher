all: test cli-build server-build

test:
	go test -v github.com/joelbirchler/label-watcher/internal

cli-build:
	go build -o lwcli cmd/cli/cli.go

cli-run: cli-build
	./lwcli

server-build:
	go build -o lwserver cmd/server/server.go

server-run: server-build
	./lwserver