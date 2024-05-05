build:
	docker-compose build

run:
	docker-compose up -d


migrateup:
	migrate -path ./db/migrations -database 'postgres://expense_user:zorro@0.0.0.0:5436/expenses?sslmode=disable' up

migratedown:
	migrate -path ./db/migrations -database 'postgres://expense_user:zorro@0.0.0.0:5436/expenses?sslmode=disable' down

createdb:
	docker exec -it postgres_exp createdb --username=expense_user --owner=expense_user expenses

dropdb:
	docker exec -it postgres_exp psql -U expense_user postgres
	drop database expenses;

backup:
	pg_dump -h localhost -p 5436 -U expense_user -d expenses > backup_db_expenses.sql
