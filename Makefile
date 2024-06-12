.PHONY: compile

PROTO_FILES = $(shell find ./investAPI/src/docs/contracts -name '*.proto')
PROTO_NAMES = $(basename $(PROTO_FILES))

GOPATH ?= $(shell go env GOPATH)
PROTOC_GEN_GO := $(GOPATH)/bin/protoc-gen-go
PROTOC := $(shell which protoc)
# If protoc isn't on the path, set it to a target that's never up to date, so
# the install command always runs.
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
# You can add instructions for other operating systems here, or use different
# branching logic as appropriate.

# If $GOPATH/bin/protoc-gen-go does not exist, we'll run this command to install
# it.
$(PROTOC_GEN_GO):
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

$(PROTO_NAMES): %: %.proto | $(PROTOC_GEN_GO) $(PROTOC)
	protoc --plugin=$(GOPATH)/bin/protoc-gen-go --go_out=. -IinvestAPI/src/docs/contracts/ ./$<

# This is a "phony" target - an alias for the above command, so "make compile"
# still works.
compile: $(PROTO_NAMES)
	go mod tidy
