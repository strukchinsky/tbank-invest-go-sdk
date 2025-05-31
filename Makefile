.PHONY: generate update

PROTO_FILES = $(shell find ./contracts -name '*.proto')
PROTO_NAMES = $(basename $(PROTO_FILES))

GOPATH := $(shell go env GOPATH)
PROTOC_GEN_GO := $(GOPATH)/bin/protoc-gen-go
PROTOC_GEN_GO_GPRC := $(GOPATH)/bin/protoc-gen-go-grpc

PROTOC := $(shell which protoc)
# If protoc isn't on the path, set it to a target that's never up to date, so the install command always runs.
ifeq ($(PROTOC),)
    PROTOC = must-rebuild
endif

# Figure out which machine we're running on.
UNAME := $(shell uname)

$(PROTOC):
# Run the right installation command for the operating system.
ifeq ($(UNAME), Darwin)
	brew install protobuf
endif
ifeq ($(UNAME), Linux)
	sudo apt-get install protobuf-compiler
endif

contracts:
	git remote add -f investapi-upstream git@github.com:RussianInvestments/investAPI.git
	git checkout -b upstream/investapi investapi-upstream/main

	# split off subdir of tracking branch into separate branch
	git subtree split -q --squash --prefix=src/docs/contracts --annotate="[investAPI] " --rejoin -b merging/investapi

	# add separate branch as subdirectory on master.
	git checkout -
	git subtree add --prefix=contracts --squash merging/investapi

# If $GOPATH/bin/protoc-gen-go does not exist, we'll run this command to install it.
$(PROTOC_GEN_GO):
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

$(PROTOC_GEN_GO_GRPC):
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

$(PROTO_NAMES): %: %.proto | $(PROTOC_GEN_GO) $(PROTOC_GEN_GO_GRPC) $(PROTOC) investAPI
	protoc --plugin=$(GOPATH)/bin/protoc-gen-go --plugin=$(GOPATH)/bin/protoc-gen-go-grpc --go-grpc_out=. --go_out=. -Icontracts/ ./$<

generate: $(PROTO_NAMES)
	go mod tidy

update: contracts
	# switch back to tracking branch, fetch & rebase.
	git checkout upstream/investapi
	git pull investapi-upstream/main

	# update the separate branch with changes from upstream
	git subtree split -q --prefix=src/docs/contracts --annotate="[investAPI] " --rejoin -b merging/investapi

	# switch back to master and use subtree merge to update the subdirectory
	git checkout -
	git subtree merge -q --prefix=contracts --squash merging/investapi
