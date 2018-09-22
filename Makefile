PACKAGE_NAME := github.com/yusufsyaifudin/go-jwt-login-example
PROJECT_DIR := $(PWD)
CURRENT_TIME := `date +%s`

TC := $(if $(TC),^$(TC)$$,"")

install-dep:
	go get -u github.com/golang/dep/cmd/dep
	go get -u github.com/twitter/go-bindata/...
	dep ensure -v

# embed main db migration
embed-migrations:
	@cd $(PROJECT_DIR)/assets/migrations && go-bindata -pkg migrations -o $(PROJECT_DIR)/assets/migrations/migrations.go *

test:
	go test -cover ./...

# make create-migration NAME="create_users_table"
create-migration:
	@[ ! -z ${NAME} ] && echo assets/migrations/$(CURRENT_TIME)_${NAME}.up.sql "\n"assets/migrations/$(CURRENT_TIME)_${NAME}.down.sql
	@touch assets/migrations/$(CURRENT_TIME)_${NAME}.up.sql assets/migrations/$(CURRENT_TIME)_${NAME}.down.sql

create-doc:
	rm -rf doc
	# build doc
	apidoc -i internal
	# use go-bindata to embed
	cd $(PROJECT_DIR)/doc && go-bindata -pkg apidoc -o $(PROJECT_DIR)/apidoc/apidoc.go . css/* fonts/* img/* locales/* utils/* vendor/*
	rm -rf doc


build: create-doc
	rm -f out/go-jwt-login-example
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o out/go-jwt-login-example $(PACKAGE_NAME)/cmd/go-jwt-login-example