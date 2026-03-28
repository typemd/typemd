.PHONY: test build-frontend

build-frontend:
	cd web/frontend && npm ci && npm run build

test: build-frontend
	go clean -testcache
	go build ./...
	go test ./...
	go vet ./...
