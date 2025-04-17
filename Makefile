postgres:
	docker run --name postgre17 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:17.2-alpine3.21

createdb:
	docker exec -it postgre17 createdb --username=root --owner=root simple_bank

migrateup:
	migrate -path db/migrations -database "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migrations -database "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

dropdb:
	docker exec -it postgre17 dropdb simple_bank

sqlc:
	sqlc generate

.PHONY:
	postgres createdb dropdb migrateup migratedown sqlc