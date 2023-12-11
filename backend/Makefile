migrateup: 
	migrate -path db/migration -database "postgresql://postgres:postgres@34.27.66.222:5432/cusocial" -verbose up

migrateup1: 
	migrate -path db/migration -database "postgresql://postgres:postgres@localhost:5432/cusocial?sslmode=disable" -verbose up 1

migratedown: 
	migrate -path db/migration -database "postgresql://postgres:postgres@34.27.66.222:5432/cusocial?sslmode=disable" -verbose down

migratedown1: 
	migrate -path db/migration -database "postgresql://postgres:postgres@localhost:5432/cusocial?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -race -cover ./... 

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go csbackend/db/sqlc Store

.PHONY:	postgres createdb dropdb migrateup migratedown sqlc server mock migrateup1 migratedown1