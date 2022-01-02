GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
BINARY_NAME=main

all: clean compress

.PHONY: build
build:
	GOOS=linux $(GOBUILD) -o $(BINARY_NAME) -v

.PHONY: compress
compress: build
	zip function.zip $(BINARY_NAME)

.PHONY: clean
clean:
	$(GOCLEAN)
	-rm function.zip
