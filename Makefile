export PROJ_PATH=github.com/alexj212/gox



export DATE := $(shell date +%Y.%m.%d-%H%M)
export LATEST_COMMIT := $(shell git log --pretty=format:'%h' -n 1)
export BRANCH := $(shell git branch |grep -v "no branch"| grep \*|cut -d ' ' -f2)
export BUILT_ON_IP := $(shell [ $$(uname) = Linux ] && hostname -i || hostname )
export BIN_DIR=./bin
export VERSION_FILE   := version.txt
export TAG     := $(shell [ -f "$(VERSION_FILE)" ] && cat "$(VERSION_FILE)" || echo '0.5.46')
export VERMAJMIN      := $(subst ., ,$(TAG))
export VERSION        := $(word 1,$(VERMAJMIN))
export MAJOR          := $(word 2,$(VERMAJMIN))
export MINOR          := $(word 3,$(VERMAJMIN))
export NEW_MINOR      := $(shell expr "$(MINOR)" + 1)
export NEW_TAG := $(VERSION).$(MAJOR).$(NEW_MINOR)


export BUILT_ON_OS=$(shell uname -a)
ifeq ($(BRANCH),)
BRANCH := master
endif

export COMMIT_CNT := $(shell git rev-list HEAD | wc -l | sed 's/ //g' )
export BUILD_NUMBER := ${BRANCH}-${COMMIT_CNT}

export COMPILE_LDFLAGS=-s -X "main.BuildDate=${DATE}" \
                          -X "main.LatestCommit=${LATEST_COMMIT}" \
                          -X "main.BuildNumber=${BUILD_NUMBER}" \
                          -X "main.BuiltOnIp=${BUILT_ON_IP}" \
                          -X "main.BuiltOnOs=${BUILT_ON_OS}"



build_info: check_prereq ## Build the container
	@echo ''
	@echo '---------------------------------------------------------'
	@echo 'BUILT_ON_IP       $(BUILT_ON_IP)'
	@echo 'BUILT_ON_OS       $(BUILT_ON_OS)'
	@echo 'DATE              $(DATE)'
	@echo 'LATEST_COMMIT     $(LATEST_COMMIT)'
	@echo 'BRANCH            $(BRANCH)'
	@echo 'COMMIT_CNT        $(COMMIT_CNT)'
	@echo 'BUILD_NUMBER      $(BUILD_NUMBER)'
	@echo 'COMPILE_LDFLAGS   $(COMPILE_LDFLAGS)'
	@echo 'PATH              $(PATH)'
	@echo 'VERSION_FILE     $(VERSION_FILE)'
	@echo 'TAG              $(TAG)'
	@echo 'VERMAJMIN        $(VERMAJMIN)'
	@echo 'VERSION          $(VERSION)'
	@echo 'MAJOR            $(MAJOR)'
	@echo 'MINOR            $(MINOR)'
	@echo 'NEW_MINOR        $(NEW_MINOR)'
	@echo 'NEW_TAG          $(NEW_TAG)'
	@echo '---------------------------------------------------------'
	@echo ''


####################################################################################################################
##
## help for each task - https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
##
####################################################################################################################
.PHONY: help

help: ## This help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help



####################################################################################################################
##
## Code vetting tools
##
####################################################################################################################


test: ## run tests
	go test -v $(PROJ_PATH)

fmt: ## run fmt on project
	#go fmt $(PROJ_PATH)/...
	gofmt -s -d -w -l .

doc: ## launch godoc on port 6060
	godoc -http=:6060

deps: ## display deps for project
	go list -f '{{ join .Deps  "\n"}}' . |grep "/" | grep -v $(PROJ_PATH)| grep "\." | sort |uniq

lint: ## run lint on the project
	golint ./...

staticcheck: ## run staticcheck on the project
	staticcheck -ignore "$(shell cat .checkignore)" .

vet: ## run go vet on the project
	go vet .

reportcard: ## run goreportcard-cli
	goreportcard-cli -v

tools: ## install dependent tools for code analysis
	go install github.com/gordonklaus/ineffassign@latest
	go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
	go install golang.org/x/lint/golint@latest
	go install github.com/gojp/goreportcard/cmd/goreportcard-cli@latest
	go install github.com/goreleaser/goreleaser@latest





add_global_module: ##add dependency to global src
	@if [ ! -d "$(GOPATH)/src/$(MODULE_DIR)" ]; then \
		echo "Adding module $(MODULE_DIR)"; \
		cd $(GOPATH); \
		go get -u $(MODULE_DIR); \
	fi


add_global_binary: ##add dependency to global bin
	@if [ ! -f "$(GOPATH)/bin/$(BINARY)" ]; then \
		echo "adding binary for $(BINARY)"; \
		go install $(BINARY_URL); \
	fi




add_global_libs: ## add global libs
	@make --no-print-directory MODULE_DIR=github.com/golang/protobuf           	add_global_module
	@make --no-print-directory MODULE_DIR=github.com/grpc-ecosystem/grpc-gateway add_global_module
	@make --no-print-directory MODULE_DIR=github.com/mwitkow/go-proto-validators add_global_module

	@make --no-print-directory BINARY_URL=github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway BINARY=protoc-gen-grpc-gateway	add_global_binary
	@make --no-print-directory BINARY_URL=github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger	 	 BINARY=protoc-gen-swagger		add_global_binary
	@make --no-print-directory BINARY_URL=github.com/mwitkow/go-proto-validators/protoc-gen-govalidators BINARY=protoc-gen-govalidators	add_global_binary
	@make --no-print-directory BINARY_URL=github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc        BINARY=protoc-gen-doc          add_global_binary

add_prerequisites: add_global_libs ## add all prerequisites








####################################################################################################################
##
## Build of binaries
##
####################################################################################################################

all: proxy test events commandrtest ## build binaries in bin dir and run tests

binaries: proxy ## build binaries in bin dir

create_dir:
	@mkdir -p $(BIN_DIR)
	@rm -f $(BIN_DIR)/web
	@ln -s ../assets $(BIN_DIR)/web


check_prereq: create_dir



build_app: create_dir
	CGO_ENABLED=0 go build -o $(BIN_DIR)/$(BIN_NAME) -a -ldflags '$(COMPILE_LDFLAGS)' $(APP_PATH)




proxy: build_info ## build proxy binary in bin dir
	@echo "build proxy"
	@cd  _examples/proxy
	make BIN_NAME=proxy APP_PATH=. build_app
	@echo ''
	@echo ''

events: build_info ## build proxy binary in bin dir
	@echo "build proxy"
	@cd  _examples/events
	make BIN_NAME=events APP_PATH=. build_app
	@echo ''
	@echo ''
commandrtest: build_info ## build proxy binary in bin dir
	@echo "build proxy"
	@cd  _examples/commandrtest
	make BIN_NAME=commandrtest APP_PATH=. build_app
	@echo ''
	@echo ''

####################################################################################################################
##
## Cleanup of binaries
##
####################################################################################################################

clean: clean_proxy  ## clean all binaries in bin dir


clean_binary: ## clean binary in bin dir
	rm -f $(BIN_DIR)/$(BIN_NAME)

clean_proxy: ## clean proxy
	make BIN_NAME=proxy clean_binary
	@rm -rf $(BIN_DIR)



publish: ## tag & push to gitlab
	@echo "\n\n\n\n\n\nRunning git add\n"
	echo "$(NEW_TAG)" > "$(VERSION_FILE)"
	git add -A
	@echo "\n\n\n\n\n\nRunning git commit\n"
	git commit -m "latest version: v$(NEW_TAG)"

	@echo "\n\n\n\n\n\nRunning git tag\n"
	git tag  "v$(NEW_TAG)"

	@echo "\n\n\n\n\n\nRunning git push\n"
	git push -f origin "v$(NEW_TAG)"

	git push -f




upgrade:
	go get -u ./...
	go mod tidy