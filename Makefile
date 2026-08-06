###### Help ###################################################################
.DEFAULT_GOAL = help

.PHONY: help
help: ## list Makefile targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

###### Setup ##################################################################
IAAS=databricks
CSB_VERSION := $(or $(CSB_VERSION), $(shell grep 'github.com/cloudfoundry/cloud-service-broker' go.mod | grep -v replace | awk '{print $$NF}' | sed -e 's/v//'))
CSB_RELEASE_VERSION := $(CSB_VERSION)

####### broker environment variables
SECURITY_USER_NAME := $(or $(SECURITY_USER_NAME), databricks-broker)
SECURITY_USER_PASSWORD := $(or $(SECURITY_USER_PASSWORD), databricks-broker-pw)

BROKER_GO_OPTS=PORT=8080 \
				DB_TYPE=sqlite3 \
				DB_PATH=/tmp/csb-db \
				SECURITY_USER_NAME=$(SECURITY_USER_NAME) \
				SECURITY_USER_PASSWORD=$(SECURITY_USER_PASSWORD) \
				DATABRICKS_HOST='$(DATABRICKS_HOST)' \
				DATABRICKS_TOKEN=$(DATABRICKS_TOKEN) \
				PAK_BUILD_CACHE_PATH=$(PAK_BUILD_CACHE_PATH) \
 				GSB_SERVICE_CSB_DATABRICKS_WORKSPACE_PLANS='$(GSB_SERVICE_CSB_DATABRICKS_WORKSPACE_PLANS)' \
				GSB_SERVICE_CSB_DATABRICKS_MODEL_SERVING_PLANS='$(GSB_SERVICE_CSB_DATABRICKS_MODEL_SERVING_PLANS)' \
 				GSB_COMPATIBILITY_ENABLE_BETA_SERVICES=$(GSB_COMPATIBILITY_ENABLE_BETA_SERVICES)

PAK_PATH := $(CURDIR)
RUN_CSB=$(BROKER_GO_OPTS) go run github.com/cloudfoundry/cloud-service-broker/v2
LDFLAGS="-X github.com/cloudfoundry/cloud-service-broker/v2/utils.Version=$(CSB_VERSION)"
GET_CSB="env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags $(LDFLAGS) github.com/cloudfoundry/cloud-service-broker/v2"

###### Targets ################################################################

.PHONY: build
build: cloud-service-broker $(IAAS)-services-*.brokerpak ## build brokerpak

$(IAAS)-services-*.brokerpak: *.yml terraform/*/*/*.tf | $(PAK_BUILD_CACHE_PATH)
	$(RUN_CSB) pak build

.PHONY: run
run: databricks_host databricks_token ## start CSB
	$(RUN_CSB) pak build --target current
	$(RUN_CSB) serve

.PHONY: docs
docs: build brokerpak-user-docs.md ## build docs

brokerpak-user-docs.md: *.yml
	$(RUN_CSB) pak docs $(PAK_PATH)/$(shell ls *.brokerpak) > $@ # GO

.PHONY: examples
examples: ## display available examples
	 $(RUN_CSB) examples

.PHONY: run-examples
run-examples: ## run examples tests, set service_name and/or example_name
	$(RUN_CSB) run-examples --service-name="$(service_name)" --example-name="$(example_name)"

.PHONY: test ## run the tests
test: lint run-integration-tests

.PHONY: run-integration-tests
run-integration-tests: ## run integration tests for this brokerpak
	cd ./integration-tests && go tool ginkgo -r .

.PHONY: run-terraform-tests
run-terraform-tests: ## run terraform tests for this brokerpak
	cd ./terraform-tests && go tool ginkgo -r .

.PHONY: info
info: build ## use the CSB to parse the buildpak and print out contents and versions
	$(RUN_CSB) pak info $(PAK_PATH)/$(shell ls *.brokerpak)

.PHONY: validate
validate: build ## use the CSB to validate the buildpak
	$(RUN_CSB) pak validate $(PAK_PATH)/$(shell ls *.brokerpak)

# fetching bits for cf push broker
.PHONY: cloud-service-broker
cloud-service-broker: go.mod ## build or fetch CSB binary
	"$(GET_CSB)"

APP_NAME := $(or $(APP_NAME), cloud-service-broker-databricks)
DB_TLS := $(or $(DB_TLS), skip-verify)


.PHONY: push-broker
push-broker: cloud-service-broker build databricks_host databricks_token ## push the broker to targeted Cloud Foundry
	MANIFEST=cf-manifest.yml APP_NAME=$(APP_NAME) DB_TLS=$(DB_TLS) ./scripts/push-broker.sh

.PHONY: databricks_host
databricks_host:
ifndef DATABRICKS_HOST
	$(error variable DATABRICKS_HOST not defined)
endif

.PHONY: databricks_token
databricks_token:
ifndef DATABRICKS_TOKEN
	$(error variable DATABRICKS_TOKEN not defined)
endif

.PHONY: clean
clean: ## clean up build artifacts
	- rm -f $(IAAS)-services-*.brokerpak
	- rm -f ./cloud-service-broker
	- rm -f ./brokerpak-user-docs.md

$(PAK_BUILD_CACHE_PATH):
	@echo "Folder $(PAK_BUILD_CACHE_PATH) does not exist. Creating it..."
	mkdir -p $@
	
.PHONY: latest-csb
latest-csb: ## point to the very latest CSB on GitHub
	go get -d github.com/cloudfoundry/cloud-service-broker@main
	go mod tidy

.PHONY: local-csb
local-csb: ## point to a local CSB repo
	echo "replace \"github.com/cloudfoundry/cloud-service-broker/v2\" => \"$$PWD/../cloud-service-broker\"" >>go.mod
	go mod tidy

.PHONY: lint
lint: checkgoformat checkgoimports checktfformat vet staticcheck ## checks format, imports and vet

checktfformat: ## checks that Terraform HCL is formatted correctly
	@@if [ "$$(terraform fmt -recursive --check)" ]; then \
		echo "terraform fmt check failed: run 'make format'"; \
		exit 1; \
	fi

checkgoformat: ## checks that the Go code is formatted correctly
	@@if [ -n "$$(gofmt -s -e -l -d .)" ]; then       \
		echo "gofmt check failed: run 'make format'"; \
		exit 1;                                       \
	fi

checkgoimports: ## checks that Go imports are formatted correctly
	@@if [ -n "$$(go tool goimports -l -d .)" ]; then \
		echo "goimports check failed: run 'make format'";                      \
		exit 1;                                                                \
	fi

vet: ## Runs go vet
	go vet ./...

staticcheck: ## Runs staticcheck
	go tool staticcheck ./...

.PHONY: format
format: ## format the source
	gofmt -s -e -l -w .
	go tool goimports -l -w .
	terraform fmt --recursive

