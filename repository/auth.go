package repository

import (
	"context"
	"crypto/sha256"
	dto "expenses/dto_expenses"
	"fmt"
	"log"
)

const (
	salt = "uiwm67#hdnl4*"
)

func (r *ExpenseRepo) CreateUser(ctx context.Context, user *dto.User) (int, error) {
	var userId int
	user.Pass = hashPassword(user.Pass)
	query := fmt.Sprintf("INSERT INTO %s (name,surname,login,pass,email) values ($1,$2,$3,$4,$5) RETURNING id", "users")
	err := r.conn.QueryRow(ctx, query, user.Name, user.Surname, user.Login, user.Pass, user.Email).Scan(&userId)
	if err != nil {
		log.Printf("error create user:%s", err.Error())
		return -1, err
	}
	return userId, nil
}

func hashPassword(pass string) string {
	hash := sha256.New()
	hash.Write([]byte(pass))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))

}
