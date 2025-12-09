# run service
.PHONY: run
run:
	go run ./...

### SQLC
.PHONY: gen-sqlc-dishes
gen-sqlc-dishes:
	sqlc -f ./schema/sqlc_configs/sqlc-dishes.yaml generate

.PHONY: gen-sqlc
gen-sqlc: gen-sqlc-dishes
