SOURCES := $(shell find . -name '*.go')

lambda/handler: go.mod go.sum $(SOURCES)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags '-extldflags "-static"' -o lambda/handler

lambda/bundle.zip: lambda/handler
	cat lambda/handler | zip bundle.zip - -i handler

build: lambda/handler
bundle: lambda/bundle.zip
