.PHONY: build-images deps run clean

CHROME_IMG = koneko-chrome
FIREFOX_IMG = koneko-firefox
# Fetch the host's CPU architecture (amd64, arm64, etc.)
ARCH := $(shell go env GOARCH)

deps:
	@echo "Installing server dependencies..."
	cd server && go mod tidy

compile-agent:
	@echo "Compiling agent for Linux on $(ARCH) architecture..."
	cd containers && GOOS=linux GOARCH=$(ARCH) go build -o koneko-agent ../agent/agent.go

build-images:compile-agent
	@echo "Building Chrome image..."
	docker build -t $(CHROME_IMG) -f containers/Dockerfile.chrome .
	@echo "Building Firefox image..."
	docker build -t $(FIREFOX_IMG) -f containers/Dockerfile.firefox .

run: deps build-images
	@echo "Starting Koneko server..."
	cd server && go run .

clean:
	@echo "Removing Koneko images..."
	docker rmi $(CHROME_IMG) $(FIREFOX_IMG) || true