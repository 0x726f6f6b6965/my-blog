PROJECTNAME := $(shell basename "$(PWD)")
include .env
export $(shell sed 's/=.*//' .env)


# db-migrate-script-gen - generate database migration files based on scripts
.PHONY: db-migrate-script-gen
db-migrate-script-gen:
	@cd db/migrations/scripts && go-bindata -pkg migrations -o ../scripts.go .