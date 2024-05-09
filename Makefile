DB_URI=postgres://expense_user:zorro@0.0.0.0:5436/expenses?sslmode=disable
DB_CONTAINER=postgres://expense_user:zorro@db:5432/expenses

build:
	docker-compose build

run:
	docker-compose up -d

migrateup:
	migrate -path ./db/schema -database "$(DB_URI)" up

migratedown:
	migrate -path ./db/schema -database "$(DB_URI)" down

backup:
	pg_dump -h localhost -p 5436 -U expense_user -d expenses > backup_db_expenses.sql

psql:
	docker exec -it postgres_exp psql "$(DB_CONTAINER)"