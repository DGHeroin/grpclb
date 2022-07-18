pwd=$(shell pwd)

$(shell mkdir -p release)
all: srv cli

srv:
	cd $(pwd)/cmd/srv/;\
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(pwd)/release/$@;

cli:
	cd $(pwd)/cmd/cli/;\
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(pwd)/release/$@;