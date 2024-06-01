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

gen-contract: solc abigen

solc: 
	@solc --bin --abi contracts/$(I).sol -o bin/contracts/ --overwrite

solc-clean:
	@rm -rf .bin/contracts/

abigen:
	@abigen --bin=bin/contracts/$(I).bin --abi=bin/contracts/$(I).abi --pkg=contract --out=contracts/gen/$(I).go
