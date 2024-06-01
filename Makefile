run: build
	@./bin/main

build:
	@templ generate components
	@tailwindcss -i view/styles.css -o assets/styles.css -m
	@go build -o bin/main main.go 