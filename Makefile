default: install

# moved out of deps to decrease build time
build-deps:
	GO111MODULE=off go get -u golang.org/x/lint/golint
	GO111MODULE=off go get -u oss.indeed.com/go/go-groups

fmt:
	go-groups -w .
	gofmt -s -w .

deps:
	go mod download

test:
	go vet ./...
	golint -set_exit_status ./...
	go test -v ./...

install:
	go install

deploy:
	mkdir -p bin
	CGO_ENABLED=0 gox -ldflags='-w -s -extldflags "-static"' -os="windows darwin" -arch="amd64 386" -output="bin/{{.Dir}}_{{.OS}}_{{.Arch}}"
	CGO_ENABLED=0 gox -ldflags='-w -s -extldflags "-static"' -os="linux" -arch="amd64 386 arm arm64" -output="bin/{{.Dir}}_{{.OS}}_{{.Arch}}"

docker:
	docker build -t mjpitz/highlander-proxy:latest .
