build: 
	@go build -o bin/bank-api

run:
	@./bin/bank-api

test: 
	@go test -v ./..