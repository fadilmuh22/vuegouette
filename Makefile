default: build-server build-client

BINARY_NAME=main

clean:
	rm -f tmp/$(BINARY_NAME)*

build-debug: clean
	CGO_ENABLED=0 go build -gcflags=all="-N -l" -o tmp/$(BINARY_NAME)-debug cmd/main.go

build-server:
	go build -o tmp/$(BINARY_NAME) cmd/main.go

build-client:
	cd web && npm i && npm run build

dev:
	docker compose up --watch