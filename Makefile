STACK_NAME ?= server-example
FUNCTIONS := products products-stream
REGION := eu-central-1

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


invoke: ## Invoke the GetProductsFunction
	@sam local invoke --env-vars env-vars.json GetProductsFunction


invoke-put: ## Invoke the PutProductFunction
	@sam local invoke --env-vars env-vars.json --event functions/put-product/event.json PutProductFunction


invoke-get: ## Invoke the GetProductFunction
	@sam local invoke --env-vars env-vars.json --event functions/get-product/event.json GetProductFunction


invoke-delete: ## Invoke the DeleteProductFunction
	@sam local invoke --env-vars env-vars.json --event functions/delete-product/event.json DeleteProductFunction


invoke-stream: ## Invoke the DDBStreamsFunction
	@sam local invoke --env-vars env-vars.json --event functions/products-stream/event.json DDBStreamsFunction


clean: ## Clean up
	@rm $(foreach function,${FUNCTIONS}, functions/${function}/bootstrap)


deploy-gateway-v2:
	if [ -f samconfig.toml ]; \
		then sam deploy --stack-name ${STACK_NAME}; \
		else sam deploy -g --stack-name ${STACK_NAME}; \
	fi


tests-unit:
	@go test -v -tags=unit -bench=. -benchmem -cover ./...


tests-integ:
	API_URL=$$(aws cloudformation describe-stacks --stack-name $(STACK_NAME) \
	  --region $(REGION) \
		--query 'Stacks[0].Outputs[?OutputKey==`ApiUrl`].OutputValue' \
		--output text) go test -v -tags=integration ./...


tests-load:
	API_URL=$$(aws cloudformation describe-stacks --stack-name $(STACK_NAME) \
	  --region $(REGION) \
		--query 'Stacks[0].Outputs[?OutputKey==`ApiUrl`].OutputValue' \
		--output text) artillery run load-testing/load-test.yml


fmt: ## Run go fmt against code
	@go fmt ./...
.PHONY: fmt

lint: ## Run the linter 
	@golangci-lint run
.PHONY: lint
