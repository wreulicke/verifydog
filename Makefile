ARCH := 386 amd64
OS := linux darwin windows

VERSION=$(shell git describe --tags 2>/dev/null)

VERSION_FLAG=\"main.version=$(shell git describe --tags 2>/dev/null)\"

status:
	dep status

install: 
	dep ensure

update:
	dep ensure -update

setup-for-ci:
	go get github.com/mitchellh/gox
	wget https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64 -O /usr/local/bin/dep
	chmod +x /usr/local/bin/dep
	
format:
	go fmt ./...

test: 
	go test ./...
	
build: 
	go build -ldflags "-X $(VERSION_FLAG)" -o ./dist/verifydog .

clean:
	rm -rf dist/

build-all: 
	gox -os="$(OS)" -arch="$(ARCH)" -ldflags "-X $(VERSION_FLAG)" -output "dist/verifydog_{{.OS}}_{{.Arch}}"
	
release:
	go get github.com/tcnksm/ghr
	@ghr -u $(CIRCLE_PROJECT_USERNAME) -r $(CIRCLE_PROJECT_REPONAME) $(VERSION) dist/
	