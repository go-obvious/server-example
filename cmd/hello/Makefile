project_root := $(shell git rev-parse --show-toplevel)
mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
app_name := $(notdir $(patsubst %/,%,$(dir $(mkfile_path))))

default: menu
.PHONY: default

menu: ## prints out the menu of command options
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(lastword $(MAKEFILE_LIST)) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
.PHONY: menu

build: ## build target application
	@bash $(project_root)/build/scripts/build.sh $(app_name)
.PHONY: build

clean: # dummy target, handled by top level makefile
.PHONY: clean

fmt:  # dummy target, handled by top level makefile
.PHONY: fmt

inspect:  # dummy target, handled by top level makefile
.PHONY: inspect

test:  # dummy target, handled by top level makefile
.PHONY: test
