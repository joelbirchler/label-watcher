build: cli-build

cli-build:
	go build -o lwcli cmd/cli/cli.go

cli-run: cli-build
	./lwcli