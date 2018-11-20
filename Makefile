.PHONY: build test clean prepare update

GO=CGO_ENABLED=0 go

MICROSERVICES=cmd/device-mqtt

.PHONY: $(MICROSERVICES)

VERSION=$(shell cat ./VERSION)

GOFLAGS=-ldflags "-X github.com/edgexfoundry/device-mqtt-go.Version=$(VERSION)"

build: $(MICROSERVICES)
	go build ./...

cmd/device-mqtt:
	$(GO) build $(GOFLAGS) -o $@ ./cmd

test:
	go test ./... -cover

clean:
	rm -f $(MICROSERVICES)

prepare:
	glide install

update:
	glide update

run:
	cd bin && ./edgex-launch.sh