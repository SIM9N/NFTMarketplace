run: build
	@./bin/main

build:
	@templ generate components
	@go build -o bin/main main.go 