


build-linux:
	GOOS=linux GOARCH=amd64 go build -o osfacts ./cmd

build-windows:
	GOOS=windows GOARCH=amd64 go build -o osfacts.exe ./cmd

build-darwin:
	GOOS=darwin GOARCH=amd64 go build -o osfacts ./cmd

build-aix:
	GOOS=aix GOARCH=ppc64 go build -o osfacts ./cmd

build-solaris:
	GOOS=solaris GOARCH=amd64 go build -o osfacts ./cmd


