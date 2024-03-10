build: 
	@go build -o bin/bank-api

run:
	@./bin/bank-api

dev: 	
	@go build -o bin/bank-api
	@./bin/bank-api

test: 
	@go test -v ./..