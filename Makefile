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

certs:
	openssl req -x509 -nodes \
		-days 365 \
		-newkey rsa:2048 \
		-subj "/C=US/ST=Oregon/L=Eugene/CN=localhost" \
		-keyout tls-key.pem \
		-out tls-cert.pem