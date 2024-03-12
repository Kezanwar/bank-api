build: 
	@go build -o bin/bank-api

run:
	@./bin/bank-api

dev: 	
	@go build -o bin/bank-api
	@./bin/bank-api

start-db: 
	@docker start bank-api-db

test: 
	@go test -v ./..