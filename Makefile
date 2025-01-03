
postgres:
	docker run --name postgres_blog -p 54332:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=aswedD4321 -d postgres:17-alpine

createdb:
	docker exec -it postgres_blog createdb --username=root --owner=root blog

dropdb:
	docker exec -it postgres_blog dropdb blog

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server
