OWNER=sjeandeaux
REPO=access-log-monitor
SRC_DIR=github.com/$(OWNER)/$(REPO)

BUILD_TIME=$(shell date +%Y-%m-%dT%H:%M:%S%z)
GIT_COMMIT?=$(shell git rev-parse --short HEAD 2> /dev/null || echo "UNKNOWN")
GIT_DIRTY?=$(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)
GIT_DESCRIBE?=$(shell git describe --tags --always 2> /dev/null || echo "UNKNOWN")
BUILD_VERSION?=$(GIT_DESCRIBE)
BUILD_TIME?=$(shell date +"%Y-%m-%dT%H:%M:%S")

LDFLAGS=-ldflags "\
          -X $(SRC_DIR)/pkg/information.Version=$(BUILD_VERSION) \
          -X $(SRC_DIR)/pkg/information.BuildTime=$(BUILD_TIME) \
          -X $(SRC_DIR)/pkg/information.GitCommit=$(GIT_COMMIT) \
          -X $(SRC_DIR)/pkg/information.GitDirty=$(GIT_DIRTY) \
          -X $(SRC_DIR)/pkg/information.GitDescribe=$(GIT_DESCRIBE)"

PKGGOFILES=$(shell go list ./... | grep -v todo-grpc)

# build in os and arch and associate to a tag
define build-and-associate
	GOOS=$(1) GOARCH=$(2) go build $(LDFLAGS) -o ./target/$(1)-$(2)-${REPO} ./${REPO}/main.go
	GOOS=$(1) GOARCH=$(2) associator $(3) -name $(1)-$(2)-${REPO} -label $(1)-$(2)-${REPO} -content-type application/binary -owner $(OWNER) -repo $(REPO) -tag $(BUILD_VERSION)  -file ./target/$(1)-$(2)-${REPO}
endef

# https://gist.github.com/sjeandeaux/e804578f9fd68d7ba2a5d695bf14f0bc
help: ## prints help.
	@grep -hE '^[a-zA-Z_-]+.*?:.*?## .*$$' ${MAKEFILE_LIST} | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: tools
tools: ## download tools
	go get -u github.com/client9/misspell/cmd/misspell
	go get -u golang.org/x/lint/golint
	go get -u github.com/fzipp/gocyclo
	go get -u gotest.tools/gotestsum
	go get golang.org/x/tools/cmd/cover
	go get github.com/mattn/goveralls
	go get github.com/sjeandeaux/toolators/cmd/associator


.PHONY: dependencies
dependencies: ## Download the dependencies
	go mod download

.PHONY: build
build: 	##Build the binary ./target/access-log-monitor
	mkdir -p ./target
	CGO_ENABLED=0 go build $(LDFLAGS) -installsuffix 'static' -o ./target/access-log-monitor ./access-log-monitor/main.go

.PHONY: gocyclo
gocyclo: ## check cyclomatic
	@gocyclo .

.PHONY: fmt
fmt: ## go fmt
	@go fmt $(PKGGOFILES)

.PHONY: misspell
misspell: ## gmisspell packages
	@misspell $(PKGGOFILES)

.PHONY: vet
vet: ## go vet on packages
	@go vet $(PKGGOFILES)

.PHONY: lint
lint: ## go lint on packages
	@golint -set_exit_status=true ./...

.PHONY: test
test: clean fmt vet ## test
	gotestsum --junitfile target/test-results/unit-tests.xml -- --short -cpu=2 -p=2 -coverprofile=target/coverage.txt -covermode=atomic -v $(LDFLAGS) $(PKGGOFILES)

.PHONY: it-test
it-test: clean fmt vet ## integration test
	gotestsum --junitfile target/test-results/it-tests.xml  -- -cpu=2 -p=2 -coverprofile=target/coverage.txt -covermode=atomic -v $(LDFLAGS) $(PKGGOFILES)

cover-html: it-test ## show the coverage in HTML page
	go tool cover -html=target/coverage.txt

clean: ## clean the target folder
	@rm -fr target
	@mkdir -p target/test-results

docker-build: ## build the docker image
	docker build --build-arg VCS_REF=$(GIT_COMMIT) --build-arg BUILD_VERSION=$(BUILD_VERSION) --build-arg BUILD_DATE=$(BUILD_TIME) --tag $(OWNER)/$(REPO):$(GIT_DESCRIBE) .

docker-push: ## push the docker image
	docker push $(OWNER)/$(REPO):$(GIT_DESCRIBE)

publish: ## publish the application in tag in progress (TODO move in circleci)
	$(call build-and-associate,linux,amd64,-create)
	$(call build-and-associate,darwin,amd64)
	$(call build-and-associate,windows,amd64)

ui-test: ## It runs the docker-compose. The flog container generates log and access-log-monitor is launched
	docker-compose up -d flog
	docker-compose run access-log-monitor