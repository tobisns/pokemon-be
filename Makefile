include .env

start:
	docker compose up --build --force-recreate

stop:
	docker compose down

migrate-up:
	migrate -path ./pkg/db/migrations/ -database "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${HOSTNAME}:${HOST_DB_PORT}/${POSTGRES_DB}?sslmode=disable" -verbose up

migrate-down:
	migrate -path ./pkg/db/migrations/ -database "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${HOSTNAME}:${HOST_DB_PORT}/${POSTGRES_DB}?sslmode=disable" -verbose down