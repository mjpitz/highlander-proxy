PROJECT ?= mjpitz--highlander-proxy
IMAGE_NAME ?= mjpitz/highlander-proxy
IMAGE_TAG ?= latest
BUILDX_OPTS ?=

default: install

# moved out of deps to decrease build time
build-deps:
	GO111MODULE=off go get honnef.co/go/tools/cmd/staticcheck
	GO111MODULE=off go get -u github.com/mitchellh/gox
	GO111MODULE=off go get -u oss.indeed.com/go/go-groups

fmt:
	go-groups -w .
	gofmt -s -w .

deps:
	go mod download

lint:
	staticcheck ./...

test:
	go vet ./...
	go test -v ./...

build:
	CGO_ENABLED=0 go build -ldflags='-w -s -extldflags "-static"' .

deploy:
	mkdir -p bin
	CGO_ENABLED=0 gox -ldflags='-w -s -extldflags "-static"' -os="windows darwin" -arch="amd64 386" -output="bin/{{.Dir}}_{{.OS}}_{{.Arch}}"
	CGO_ENABLED=0 gox -ldflags='-w -s -extldflags "-static"' -os="linux" -arch="amd64 386 arm arm64" -output="bin/{{.Dir}}_{{.OS}}_{{.Arch}}"

docker:
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) .

dockerx:
	docker buildx create --name $(PROJECT) &>/dev/null || echo "$(PROJECT) project exists"
	docker buildx use $(PROJECT)
	docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 -t $(IMAGE_NAME):$(IMAGE_TAG) . $(BUILDX_OPTS)

publish:
	git checkout $(VERSION)
	rm -rf bin/
	make deploy
	BUILDX_OPTS=$(BUILDX_OPTS) IMAGE_TAG=$(VERSION) make dockerx
	BUILDX_OPTS=$(BUILDX_OPTS) IMAGE_TAG=latest make dockerx
