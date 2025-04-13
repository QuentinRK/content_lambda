.PHONY: build run test clean docker-build docker-run

build:
	go build -o bin/main ./cmd/lambda

run: build
	./bin/main

test:
	go test ./...

clean:
	rm -rf bin/
	go clean

docker-build:
	docker-compose build

docker-run:
	docker-compose up

docker-down:
	docker-compose down

# Initialize DynamoDB tables
init-db:
	aws dynamodb create-table \
		--endpoint-url http://localhost:8000 \
		--table-name Content \
		--attribute-definitions AttributeName=ID,AttributeType=S \
		--key-schema AttributeName=ID,KeyType=HASH \
		--provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 