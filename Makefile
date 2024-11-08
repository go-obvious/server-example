STACK_NAME ?= server-example
FUNCTIONS := products products-stream
REGION := us-east-1

GO := go

default: menu
.PHONY: default

menu:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(lastword $(MAKEFILE_LIST)) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
.PHONY: menu

.PHONY: build tests-unit tests-integ tests-load deploy-gateway-v2 invoke invoke-put invoke-get invoke-delete invoke-stream clean lint

ci: build tests-unit

build: ## Build the functions
		${MAKE} ${MAKEOPTS} $(foreach function,${FUNCTIONS}, build-${function})


build-%:
		cd functions/$* && GOOS=linux GOARCH=arm64 CGO_ENABLED=0 ${GO} build -o bootstrap


clean: ## Clean up
	@rm $(foreach function,${FUNCTIONS}, functions/${function}/bootstrap)


deploy-gateway-v2: ## Deploy the API Gateway v2
	if [ -f samconfig.toml ]; \
		then samlocal deploy --stack-name ${STACK_NAME}; \
		else samlocal deploy -g --stack-name ${STACK_NAME}; \
	fi

destroy-gateway-v2: ## Destroy the API Gateway v2
	samlocal delete --stack-name ${STACK_NAME}


tests-unit: ## Run the unit tests
	@go test -v -tags=unit -bench=. -benchmem -cover ./...


tests-integ: ## Run the integration tests
	API_URL=$$(awslocal cloudformation describe-stacks --stack-name $(STACK_NAME) \
		--region $(REGION) \
		--query 'Stacks[0].Outputs[?OutputKey==`ApiUrl`].OutputValue' \
		--output text) go test -v -tags=integration ./...


tests-load: ## Run the load tests
	API_URL=$$(awslocal cloudformation describe-stacks --stack-name $(STACK_NAME) \
		--region $(REGION) \
		--query 'Stacks[0].Outputs[?OutputKey==`ApiUrl`].OutputValue' \
		--output text) artillery run load-testing/load-test.yml

echo: ## Output the API URL
	awslocal cloudformation describe-stacks --stack-name $(STACK_NAME) \
		--region $(REGION) \
		--query 'Stacks[0].Outputs[?OutputKey==`ApiUrl`].OutputValue' \
		--output text


fmt: ## Run go fmt against code
	@go fmt ./...
.PHONY: fmt

lint: ## Run the linter 
	@golangci-lint run
.PHONY: lint
