.PHONY: build test clean prepare update

GO=CGO_ENABLED=0 go
GOFLAGS=-ldflags

MICROSERVICES=cmd/device-mqtt

.PHONY: $(MICROSERVICES)

build: $(MICROSERVICES)
	go build ./...

cmd/device-mqtt:
	$(GO) build -o $@ ./cmd

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