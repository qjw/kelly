REPO = github.com/qjw/kelly
BINARIES:=helloworld response route

all: $(BINARIES)
build: $(BINARIES)

$(BINARIES):
	mkdir -p build
	go build -mod=vendor -o build "./examples/$@"

.PHONY: mod
mod:
	go mod tidy && go mod vendor
