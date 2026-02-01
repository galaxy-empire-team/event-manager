.PHONY: local-run
local-run:
	bash -c 'set -a; . .env; set +a; go run cmd/manager/main.go'

.PHONY: lint
lint:
	golangci-lint run -v
