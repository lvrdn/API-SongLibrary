run:
	go run cmd/main.go

db:
	psql -p5432 -Uroot -dsonglibrary

container:
	docker run --name postgres -e POSTGRES_USER=root -e POSTGRES_PASSWORD=1234 -p 5432:5432 -d postgres

createdb:
	docker exec -it postgres createdb --username=root --owner=root songlibrary

dropdb:
	docker exec -it postgres dropdb songlibrary

migrateup:
	migrate -path ./migration -database "postgresql://root:1234@localhost:5432/songlibrary?sslmode=disable" -verbose up

migratedown:
	migrate -path ./migration -database "postgresql://root:1234@localhost:5432/songlibrary?sslmode=disable" -verbose down

swag:
	swag init -g cmd/main.go
