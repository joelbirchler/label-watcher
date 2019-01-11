all: test cli-build

test:
	go test -v github.com/joelbirchler/label-watcher/internal

cli-build:
	go build -o lwcli cmd/cli/cli.go

cli-run: cli-build
	./lwcli