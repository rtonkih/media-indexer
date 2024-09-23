test:
	go test ./... -v

migration:
	go run ./tools/migrate/migrate.go

populate:
	go run ./tools/populate/populate.go

run:
	MODE=prod GIN_MODE=release docker-compose up --build

run-dev:
	MODE=dev docker-compose up
