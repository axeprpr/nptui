BINARY_NAME=nptui
GIT_COMMITS := $(shell git rev-list --count HEAD 2>/dev/null || echo "0")
VERSION=1.0.$(GIT_COMMITS)
BUILD_DIR=build
PACKAGE_DIR=package

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build flags
LDFLAGS=-ldflags "-s -w"

.PHONY: all build clean test deps build-amd64 build-arm64 deb-amd64 deb-arm64 deb-all install

all: clean deps build

deps:
	$(GOMOD) download
	$(GOMOD) tidy

build:
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .

build-amd64:
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-amd64 .

build-arm64:
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-arm64 .

clean:
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -rf $(PACKAGE_DIR)

test:
	$(GOTEST) -v ./...

# Install to /usr/local/bin
install: build
	install -m 755 $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/

# Create Debian package for amd64
deb-amd64: build-amd64
	@echo "Building Debian package for amd64..."
	@mkdir -p $(PACKAGE_DIR)/$(BINARY_NAME)-$(VERSION)-amd64/DEBIAN
	@mkdir -p $(PACKAGE_DIR)/$(BINARY_NAME)-$(VERSION)-amd64/usr/bin
	@mkdir -p $(PACKAGE_DIR)/$(BINARY_NAME)-$(VERSION)-amd64/usr/share/doc/$(BINARY_NAME)
	
	@cp $(BUILD_DIR)/$(BINARY_NAME)-amd64 $(PACKAGE_DIR)/$(BINARY_NAME)-$(VERSION)-amd64/usr/bin/$(BINARY_NAME)
	@chmod 755 $(PACKAGE_DIR)/$(BINARY_NAME)-$(VERSION)-amd64/usr/bin/$(BINARY_NAME)
	
	@cp debian/control-amd64 $(PACKAGE_DIR)/$(BINARY_NAME)-$(VERSION)-amd64/DEBIAN/control
	@cp debian/postinst $(PACKAGE_DIR)/$(BINARY_NAME)-$(VERSION)-amd64/DEBIAN/postinst
	@cp debian/postrm $(PACKAGE_DIR)/$(BINARY_NAME)-$(VERSION)-amd64/DEBIAN/postrm
	@chmod 755 $(PACKAGE_DIR)/$(BINARY_NAME)-$(VERSION)-amd64/DEBIAN/postinst
	@chmod 755 $(PACKAGE_DIR)/$(BINARY_NAME)-$(VERSION)-amd64/DEBIAN/postrm
	
	@cp README.md $(PACKAGE_DIR)/$(BINARY_NAME)-$(VERSION)-amd64/usr/share/doc/$(BINARY_NAME)/
	@cp debian/copyright $(PACKAGE_DIR)/$(BINARY_NAME)-$(VERSION)-amd64/usr/share/doc/$(BINARY_NAME)/
	
	@dpkg-deb --build $(PACKAGE_DIR)/$(BINARY_NAME)-$(VERSION)-amd64
	@mv $(PACKAGE_DIR)/$(BINARY_NAME)-$(VERSION)-amd64.deb $(BUILD_DIR)/
	@echo "Package created: $(BUILD_DIR)/$(BINARY_NAME)-$(VERSION)-amd64.deb"

# Create Debian package for arm64
deb-arm64: build-arm64
	@echo "Building Debian package for arm64..."
	@mkdir -p $(PACKAGE_DIR)/$(BINARY_NAME)-$(VERSION)-arm64/DEBIAN
	@mkdir -p $(PACKAGE_DIR)/$(BINARY_NAME)-$(VERSION)-arm64/usr/bin
	@mkdir -p $(PACKAGE_DIR)/$(BINARY_NAME)-$(VERSION)-arm64/usr/share/doc/$(BINARY_NAME)
	
	@cp $(BUILD_DIR)/$(BINARY_NAME)-arm64 $(PACKAGE_DIR)/$(BINARY_NAME)-$(VERSION)-arm64/usr/bin/$(BINARY_NAME)
	@chmod 755 $(PACKAGE_DIR)/$(BINARY_NAME)-$(VERSION)-arm64/usr/bin/$(BINARY_NAME)
	
	@cp debian/control-arm64 $(PACKAGE_DIR)/$(BINARY_NAME)-$(VERSION)-arm64/DEBIAN/control
	@cp debian/postinst $(PACKAGE_DIR)/$(BINARY_NAME)-$(VERSION)-arm64/DEBIAN/postinst
	@cp debian/postrm $(PACKAGE_DIR)/$(BINARY_NAME)-$(VERSION)-arm64/DEBIAN/postrm
	@chmod 755 $(PACKAGE_DIR)/$(BINARY_NAME)-$(VERSION)-arm64/DEBIAN/postinst
	@chmod 755 $(PACKAGE_DIR)/$(BINARY_NAME)-$(VERSION)-arm64/DEBIAN/postrm
	
	@cp README.md $(PACKAGE_DIR)/$(BINARY_NAME)-$(VERSION)-arm64/usr/share/doc/$(BINARY_NAME)/
	@cp debian/copyright $(PACKAGE_DIR)/$(BINARY_NAME)-$(VERSION)-arm64/usr/share/doc/$(BINARY_NAME)/
	
	@dpkg-deb --build $(PACKAGE_DIR)/$(BINARY_NAME)-$(VERSION)-arm64
	@mv $(PACKAGE_DIR)/$(BINARY_NAME)-$(VERSION)-arm64.deb $(BUILD_DIR)/
	@echo "Package created: $(BUILD_DIR)/$(BINARY_NAME)-$(VERSION)-arm64.deb"

# Build all packages
deb-all: deb-amd64 deb-arm64
	@echo "All packages built successfully!"
	@ls -lh $(BUILD_DIR)/*.deb

