VERSION="$(shell head -1 VERSION.md)"
GIT_COMMIT=$(shell git rev-parse HEAD)
COMPILE_DATE=$(shell date -u +%Y%m%d.%H%M%S)
VERSION_PACKAGE_PATH="github.com/faja/tornimo-agent/cmd/agent"
BUILD_FLAGS=-X $(VERSION_PACKAGE_PATH).versionInfoDate=$(COMPILE_DATE) -X $(VERSION_PACKAGE_PATH).versionInfoCommit=$(GIT_COMMIT) -X $(VERSION_PACKAGE_PATH).versionInfoCli=$(VERSION)

build:
	go build -ldflags "$(BUILD_FLAGS)" -o tornimo-agent cmd/main.go
