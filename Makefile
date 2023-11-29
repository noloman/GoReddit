.PHONY: postgres migrate

postgres:
	docker run -ti --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=secret -e POSTGRES_DB=postgres -e POSTGRES_USER=postgres postgres

migrate:
	migrate -source file://migrations \
		-database postgres://postgres:secret@localhost/postgres?sslmode=disable up
