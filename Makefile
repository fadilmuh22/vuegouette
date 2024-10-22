default: build_server build_client

build_server:
	go build -o main cmd/main.go

build_client:
	cd web && npm i && npm run build

dev:
	docker compose up --watch