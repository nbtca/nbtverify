export PATH := $(GOPATH)/bin:$(PATH)
export GO111MODULE=on
LDFLAGS := -s -w

all: fmt build

build: nbtverify

fmt:
	go fmt ./...

fmt-more:
	gofumpt -l -w .

vet:
	go vet ./...
	
nbtverify:
	env CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -o bin/nbtverify .

clean:
	rm -f ./bin/nbtverify
