PROJECTNAME := $(shell basename "$(PWD)")
include .env
export $(shell sed 's/=.*//' .env)


# db-migrate-script-gen - generate database migration files based on scripts
.PHONY: db-migrate-script-gen
db-migrate-script-gen:
	@cd db/migrations/scripts && go-bindata -pkg migrations -o ../scripts.go .

## Dockerfile
.PHONY: gen-images
gen-images:
	@docker build --tag blog-svc:$(BLOG_VERSION) -f ./blog-service/build/Dockerfile .
	@docker build --tag search-svc:$(SEARCH_VERSION) -f ./search-service/build/Dockerfile .
	@docker build --tag user-svc:$(USER_VERSION) -f ./user-service/build/Dockerfile .
	@docker build --tag graph-svc:$(GRAPH_VERSION) -f ./graph-service/build/Dockerfile .

.PHONY: service-up
service-up:
	@docker-compose up -d
	@docker-compose --profile tools run --rm migrate up

.PHONY: service-down
service-down:
	@docker-compose --profile tools run --rm migrate down
	@docker-compose down