NAME = $(shell basename  ${PWD})

build:
	CGO_ENABLED=0 go build -o bin/${NAME}

build-windows:
	GOOS=windows GOARCH=amd64 go build -o bin/${NAME}.exe