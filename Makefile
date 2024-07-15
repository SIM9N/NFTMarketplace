run-app: build-app
	@./bin/app/main

build-app:
	@templ generate components
	@tailwindcss -i view/styles.css -o assets/styles.css -m
	@go build -o bin/app/main ./cmd/app/main.go 

run-dev:
	@air

deploy-contract: compile-contract
	@go build -o bin/deploy/deployNFT721 ./cmd/deploy/deployNFT721.go 
	@./bin/deploy/deployNFT721

compile-contract: 
	@solc --bin --abi contracts/NFT721.sol -o bin/contracts/ --overwrite --base-path ./ --include-path ./node_modules/
	@abigen --bin=bin/contracts/NFT721.bin --abi=bin/contracts/NFT721.abi --pkg=NFT721 --out=contracts/gen/NFT721.go

clean:
	@rm -rf ./contracts/gen/*
	@rm -rf ./bin

hardhat-up:
	@npx hardhat node