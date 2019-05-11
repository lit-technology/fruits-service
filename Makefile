PORT?=8000
APPLICATION:=fruits-services
PACKAGE:=github.com/philip-bui/${APPLICATION}
PACKAGES=`go list ./...`
COVERAGE:=coverage.out
POSTGRES_DOCKER:=postgres:11.1
POSTGRES_NAME:=postgres
POSTGRES_USER:=postgres
POSTGRES_PW:=postgres
POSTGRES_DB:=fruits
POSTGRES_DW_DB:=fruits_dw

help: ## Display this help screen
	grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

godoc: ## Open HTML API Documentation
	echo "localhost:${PORT}/pkg/${PACKAGE}"
	godoc -http=:${PORT}

setup: mod protos postgres db tables

mod: # Downloads dependencies
	export GOMODULES111=on
	go get ${PACKAGES}

zip: ## Zip Package for ELB
	env GOOS=linux GOARCH=amd64 go build
	mv ${APPLICATION} application
	zip application.zip application certs/* static/* config/serviceAccountKey.json config/creds.json

deploy: zip ## Deploy to ELB
	eb deploy

test: ## Run Tests
	go test -coverprofile=coverage.out ./...

benchmark: ## Run Benchmark Tests
	go test -v -bench=. ./..

coverage: test ## Open HTML Test Coverage Report
	go tool cover -html=${COVERAGE}

proto: ## Generate Protobuf files
	go install github.com/gogo/protobuf/protoc-gen-gofast
	protoc -I protos/ protos/*.proto --gofast_out=.:protos

json: ## Generate JSON files
	go install github.com/francoispqt/gojay/gojay
	gojay -s protos/ -p True -t Survey,SurveyAnswer -pkg fruits -o protos/${APPLICATION}.json.go

postgres-stop: ## Stop Postgres Docker
	docker ps -aq --filter name=${POSTGRES_NAME} | xargs docker stop

postgres: ## Run Postgres Docker
	docker run --name ${POSTGRES_NAME} -d -p 5432:5432 ${POSTGRES_DOCKER}

db: ## Create Postgres Databases
	docker exec -i ${POSTGRES_NAME} psql -U ${POSTGRES_USER} -c 'CREATE DATABASE ${POSTGRES_DB}'
	docker exec -i ${POSTGRES_NAME} psql -U ${POSTGRES_USER} -c 'CREATE DATABASE ${POSTGRES_DW_DB}'

tables: ## Insert Postgres Tables
	docker cp ./resources/ postgres:/
	for sql in `ls -A1 ./resources/fruits`; do \
		docker exec -it ${POSTGRES_NAME} psql -U ${POSTGRES_USER} -d ${POSTGRES_DB} -f resources/fruits/$$sql; \
	done;
	for sql in `ls -A1 ./resources/dw`; do \
		docker exec -it ${POSTGRES_NAME} psql -U ${POSTGRES_USER} -d ${POSTGRES_DW_DB} -f resources/dw/$$sql; \
	done;

