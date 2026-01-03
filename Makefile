# ============================================================================
# run service
# ============================================================================

.PHONY: build
build:
	go build -o calorina ./internal/cmd/*

.PHONY: run
run:
	go run ./...

# ============================================================================
# SQLC
# ============================================================================

.PHONY: gen-sqlc-dishes
gen-sqlc-dishes:
	sqlc -f ./schema/sqlc_configs/sqlc-dishes.yaml generate

.PHONY: gen-sqlc
gen-sqlc: gen-sqlc-dishes

# ============================================================================
# Linting
# ============================================================================

.PHONY: lint
lint:
	@golangci-lint run --config=.golangci.yml

.PHONY: lint-fix
lint-fix:
	@golangci-lint run --config=.golangci.yml --fix
