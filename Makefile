


build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o osfacts ./cmd

build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o osfacts.exe ./cmd

build-darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o osfacts ./cmd

build-aix:
	CGO_ENABLED=0 GOOS=aix GOARCH=ppc64 go build -o osfacts ./cmd

build-solaris:
	CGO_ENABLED=0 GOOS=solaris GOARCH=amd64 go build -o osfacts ./cmd


