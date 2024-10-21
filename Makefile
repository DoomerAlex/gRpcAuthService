run:
	go run cmd/sso/main.go --config=./config/local.yaml

migrate:
	go run ./cmd/migrator --storage-path=./storage/sso.db --migrations-path=./migrations

test-migrate:
	go run ./cmd/migrator --storage-path=./storage/sso.db --migrations-path=./tests/migrations --migrations-table=migrations_test