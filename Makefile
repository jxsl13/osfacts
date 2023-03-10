


build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ./...

build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build ./...

build-darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build ./...

build-aix:
	CGO_ENABLED=0 GOOS=aix GOARCH=ppc64 go build ./...

build-solaris:
	CGO_ENABLED=0 GOOS=solaris GOARCH=amd64 go build ./...


build-all: build-linux build-windows build-darwin build-aix build-solaris