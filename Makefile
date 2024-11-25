
#create shupir container of postgres image
postgres:
	docker run --name sh -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres

#create database in sh container
createdb:
	docker exec -it sh createdb --username=root --owner=root shupir

# delete shupir database
drop:
	docker exec -it sh dropdb shupir

migrate:
	migrate create -ext sql -dir db/migration -seq $(name)
	
migrateup:
	migrate -path db/migration -database 'postgresql://root:secret@localhost:5432/shupir?sslmode=disable' -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/shupir?sslmode=disable" -verbose down

#execute sql queries from db/query 
sqlc:
	sqlc generate

redis:
	docker run --name sh-stikers -d redis:7-alpine

redis-cli:
	docker exec -it sh-stikers redis-cli

test:
	go test -v -cover ./...

mock:
	mockgen -source=internal/service/service.go \
	-destination=internal/service/mock/mock_service.go
lint:
	sh run-linter.sh
#run application
server:
	go run ./cmd/api


.PHONY: server postgres createdb dropdb migrateup migratedown test sqlc mock lint






