build:
	docker-compose build

run:
	docker-compose up


migrateup:
	migrate -path ./db/migrations -database 'postgres://expense_user:zorro@0.0.0.0:5436/expenses?sslmode=disable' up

migratedown:
	migrate -path ./db/migrations -database 'postgres://expense_user:zorro@0.0.0.0:5436/expenses?sslmode=disable' down


