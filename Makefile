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
	curl -sSfL https://raw.githubusercontent.com/truvami/decoder/main/install.sh | sh -s -- -b $(go env GOPATH)/bin

check-json-tags:
	@echo "Checking JSON tags for camelCase format..."
	@bash -c ' \
		files=$$(find . -name "*.go" -not -path "./vendor/*"  -not -path "./pkg/solver/loracloud/*" -not -path "./internal/selfupdate/selfupdate.go"); \
		camel_case_regex="^[a-z]+([A-Za-z0-9]+)*$$"; \
		error_found=false; \
		for file in $$files; do \
			json_tags=$$(grep -o '"'"'json:"[^"]*'"'"' "$$file" | sed '"'"'s/json:"//; s/"//'"'"'); \
			for tag in $$json_tags; do \
				if [[ ! "$$tag" =~ $$camel_case_regex ]]; then \
					echo "❌ JSON tag \"$$tag\" in file \"$$file\" is not camelCase."; \
					error_found=true; \
				fi; \
			done; \
		done; \
		if [ "$$error_found" = true ]; then \
			echo "🚧 Some JSON tags are not camelCase. Please fix them."; \
			exit 1; \
		else \
			echo "✅ All JSON tags are camelCase!"; \
		fi \
	'

check-metrics:
	@echo "🔍 Checking Prometheus metrics for 'truvami_' prefix..."
	@bad_metrics=$$(grep -r --include="*.go" -E 'prometheus\.(CounterOpts|GaugeOpts|HistogramOpts|SummaryOpts)' . | cut -d: -f1 | sort -u | xargs grep -n 'Name:' | grep -v -E 'Name:.*"truvami_'); \
	if [ -n "$$bad_metrics" ]; then \
		echo "❌ ERROR: Found Prometheus metrics without 'truvami_' prefix:"; \
		echo "$$bad_metrics"; \
		exit 1; \
	else \
		echo "✅ All Prometheus metrics are correctly prefixed."; \
	fi

.PHONY: check-coverage check-json-tags check-metrics
