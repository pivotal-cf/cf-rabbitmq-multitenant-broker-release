###### Help ###################################################################

.DEFAULT_GOAL = help

.PHONY: help

help:  ## list Makefile targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

###### Tools ###############################################################

go-vet: ## runs go vet in source code
	go vet `go list ./... | grep -v vendor`

run: ## starts the broker for testing
	go run -mod=vendor main.go -configPath ${config_path}

integration_tests: ## runs integration tests in Docker (for development)
	gcloud --project cf-rabbitmq auth configure-docker
	docker run --rm -v "${PWD}:/rabbitmq-service-broker" gcr.io/cf-rabbitmq/pcf-rabbitmq-ci:stable /rabbitmq-service-broker/ci/run_integration_tests.sh

unit_tests: ## runs the unit tests
	GOFLAGS='-mod=vendor' ginkgo -r --skip-package=integrationtests

tests: unit_tests integration_tests ## runs unit and integration tests

coverage: ## runs test coverage tool
	GOFLAGS='-mod=vendor' ginkgo -r --skip-package=integrationtests --cover --coverprofile rabbitmq-service-broker.coverprofile --outputdir .
	echo 'mode: atomic' >rabbitmq-service-broker.coverage ## https://github.com/onsi/ginkgo/issues/518
	grep -v 'mode:' rabbitmq-service-broker.coverprofile >>rabbitmq-service-broker.coverage
	go tool cover -html=rabbitmq-service-broker.coverage

ci_tests: unit_tests
	./ci/run_integration_tests.sh

deps:
	go get -u github.com/maxbrunsfeld/counterfeiter
