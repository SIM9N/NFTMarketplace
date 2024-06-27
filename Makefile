run-app: build-app
	@./bin/app/main

build-app:
	@templ generate components
	@tailwindcss -i view/styles.css -o assets/styles.css -m
	@go build -o bin/app/main ./cmd/app/main.go 

run-dev:
	@air

deploy-contract: compile-contract
	@go build -o bin/deploy/$(C) ./cmd/deploy/$(C).go 
	@./bin/deploy/$(C)

compile-contract: 
	@solc --bin --abi contracts/$(C).sol -o bin/contracts/ --overwrite --base-path ./ --include-path ./node_modules/
	@abigen --bin=bin/contracts/$(C).bin --abi=bin/contracts/$(C).abi --pkg=$(C) --out=contracts/gen/$(C).go

clean:
	@rm -rf ./contracts/gen/*
	@rm -rf ./bin

hardhat-up:
	@npx hardhat node
