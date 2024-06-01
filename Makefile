run-app: build-app
	@./bin/app/main

build-app:
	@templ generate components
	@tailwindcss -i view/styles.css -o assets/styles.css -m
	@go build -o bin/app/main ./cmd/app/main.go 


run-contract: build-contract
	@./bin/contract/main

build-contract:
	@go build -o bin/contract/main ./cmd/contract/main.go 
