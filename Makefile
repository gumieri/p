MODULE=github.com/gumieri/p
GOBUILD=go build -ldflags "-X ${MODULE}/cmd.Version=${PROJECT_VERSION}"
all: deps build
deps:
	go get ./...
build:
	$(GOBUILD)
install:
	go install
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
build-macos-86:
	GOOS=darwin \
	GOARCH=386 \
	$(GOBUILD) -o release/p-Darwin-i386
build-all: build-linux-64 build-linux-86 build-linux-arm5 build-linux-arm6 build-linux-arm7 build-linux-arm8 build-macos-64 build-macos-86
