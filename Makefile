
build-linux:
	@GOOS=linux GOARCH=amd64 go build -o bin/linux/schema .

build-mac:
	@GOOS=darwin GOARCH=amd64 go build -o bin/mac/schema .