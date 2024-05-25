package repository

import (
	"context"
	dto "expenses/dto_expenses"
	"fmt"
	"log"
)

func (r *ExpenseRepo) CreateUser(ctx context.Context, user *dto.User) (int, error) {
	var userId int

	query := fmt.Sprintf("INSERT INTO %s (name,surname,login,pass,email) values ($1,$2,$3,$4,$5) RETURNING id", "users")
	err := r.conn.QueryRow(ctx, query, user.Name, user.Surname, user.Login, user.Pass, user.Email).Scan(&userId)
	if err != nil {
		log.Printf("error create user:%s", err.Error())
		return -1, err
	}
	return userId, nil
}

func (r *ExpenseRepo) GetUser(userName, hashPassword string) (*dto.User, error) {
	var user dto.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE name=$1 and pass=$2", "users")
	err := r.conn.QueryRow(context.Background(), query, userName, hashPassword).Scan(&user.Id)
	if err != nil {
		log.Println("error get user", err.Error())
		return nil, err
	}
	return &user, nil
}
