MODULE=github.com/gumieri/p
GOBUILD=go build -ldflags "-X ${MODULE}/cmd.Version=${PROJECT_VERSION}"

define RELEASE_BODY
Go to the [README](https://github.com/gumieri/p/blob/${PROJECT_VERSION}/README.md) to know how to use it.
If you are using Mac OS or Linux, you can install using the commands:
```bash
curl -L https://github.com/gumieri/p/releases/download/${PROJECT_VERSION}/p-`uname -s`-`uname -m` -o /usr/local/bin/p
chmod +x /usr/local/bin/p
```
If you already have an older version installed, just run:
```bash
p upgrade
```
endef
export RELEASE_BODY

all: deps build
deps:
	go get ./...
build:
	$(GOBUILD)
install:
	go install
release-body:
	echo "$$RELEASE_BODY" > RELEASE.md
build-linux-64:
	GOOS=linux \
	GOARCH=amd64 \
	$(GOBUILD) -o release/p-Linux-x86_64
build-linux-86:
	GOOS=linux \
	GOARCH=386 \
	$(GOBUILD) -o release/p-Linux-i386
build-linux-arm5:
	GOOS=linux \
	GOARCH=arm \
	GOARM=5 \
	$(GOBUILD) -o release/p-Linux-armv5l
build-linux-arm6:
	GOOS=linux \
	GOARCH=arm \
	GOARM=6 \
	$(GOBUILD) -o release/p-Linux-armv6l
build-linux-arm7:
	GOOS=linux \
	GOARCH=arm \
	GOARM=7 \
	$(GOBUILD) -o release/p-Linux-armv7l
build-linux-arm8:
	GOOS=linux \
	GOARCH=arm64 \
	$(GOBUILD) -o release/p-Linux-armv8l
build-macos-64:
	GOOS=darwin \
	GOARCH=amd64 \
	$(GOBUILD) -o release/p-Darwin-x86_64
build-macos-arm64:
	GOOS=darwin \
	GOARCH=arm64 \
	$(GOBUILD) -o release/p-Darwin-arm-64
build-all: build-linux-64 build-linux-86 build-linux-arm5 build-linux-arm6 build-linux-arm7 build-linux-arm8 build-macos-64 build-macos-arm64
