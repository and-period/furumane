.PHONY: help setup install fmt vet lint test build proto mockgen

GOFUMPT_VERSION := 0.5.0
GOLANGCI_VERSION := 1.54.2
MOCKGEN_VERSION := 0.2.0

LINT_PACKAGES := $(shell go list $(CURDIR)/... | grep -v -e "mock" -v -e "proto" -v -e "tmp")
TEST_PACKAGES := $(shell go list $(CURDIR)/internal/... $(CURDIR)/pkg/...)

help:
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

setup: proto install ## 初回環境構築用
	wget -O - -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v${GOLANGCI_VERSION}

install: ## 依存ライブラリのインストール
	go install mvdan.cc/gofumpt@v${GOFUMPT_VERSION}
	go install go.uber.org/mock/mockgen@v${MOCKGEN_VERSION}

clean:
	rm -rf ./proto/**/*.pb.go ./proto/**/*.pb.*.go

fmt: ## フォーマットが正しくない箇所の出力
	! gofumpt -d ./cmd ./config ./hack ./internal ./pkg | grep '^'

vet: ## コードの静的解析
	go vet $(LINT_PACKAGES)

lint: ## Linterの実行
	buf lint
	./bin/golangci-lint run -c .golangci.yaml ./...

test: ## テストの実行
	go test -v -cover -coverprofile=coverage.txt -covermode=atomic $(TEST_PACKAGES)

build: ## アプリケーションのコンパイル
	go build -o ./app ./cmd/${SERVICE}/main.go

proto: clean ## Protocol Buffersの定義からファイル生成
	buf generate

mockgen: ## ユニットテストで使用するモックの生成
	rm -rf ./mock
	go generate ./...
	./hack/generate-mocks.sh
