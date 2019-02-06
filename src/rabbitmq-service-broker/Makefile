###### Help ###################################################################

.DEFAULT_GOAL = help

.PHONY: help

help:  ## list Makefile targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

###### Tools ###############################################################

go-vet: ## runs go vet in source code
	go vet `go list ./... | grep -v vendor`

run: ## starts the broker for testing
	go run -mod=vendor cmd/rabbitmq-service-broker/main.go
