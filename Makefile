OS := $(shell uname | awk '{print tolower($$0)}')
MACHINE := $(shell uname -m)

.DEFAULT_GOAL:=help

# set default shell
SHELL=/bin/bash -o pipefail -o errexit

.PHONY: help build clean

help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

build:  ## Build this go pacakge
	@pushd example/pineapple > /dev/null; \
	go build; \
	popd > /dev/null

example: build  ## Run the "hello, world" example
	@pushd example/pineapple > /dev/null; \
	echo "This is the original code:"; \
	cat helloworld.papp; \
	echo -e "\nRun this program, and the following are results!"; \
	./pineapple helloworld.papp; \
	popd > /dev/null

clean:  ## Clean the binaries
	@rm -f example/pineapple/pineapple
