SHELL:=/bin/bash

YELLOW := \e[0;33m
RESET := \e[0;0m

GOVER := $(shell go env GOVERSION)
GOMINOR := $(shell bash -c "cut -f2 -d. <<< $(GOVER)")

define execute-if-go-122
@{ \
if [[ 22 -le $(GOMINOR) ]]; then \
	$1; \
else \
	echo -e "$(YELLOW)Skipping task as you're running Go v1.$(GOMINOR).x which is < Go 1.22, which this module requires$(RESET)"; \
fi \
}
endef

coverage:
	go test -coverprofile cover.out `go list ./...`
	go tool cover -html=cover.out

GOBIN ?= $$(go env GOPATH)/bin

.PHONY: install-go-test-coverage
install-go-test-coverage:
	go install github.com/vladopajic/go-test-coverage/v2@latest

.PHONY: check-coverage
check-coverage: install-go-test-coverage
	export ENV=test && go test ./... -coverprofile=./cover.out -covermode=atomic -coverpkg=./...
	${GOBIN}/go-test-coverage --config=./.testcoverage.yml

install:
<<<<<<< Updated upstream
	curl -sSfL https://raw.githubusercontent.com/truvami/decoder/main/install.sh | sh -s -- -b $(go env GOPATH)/bin
=======
	curl -sSfL https://raw.githubusercontent.com/truvami/decoder/main/install.sh | sh -s -- -b $(go env GOPATH)/bin
>>>>>>> Stashed changes
