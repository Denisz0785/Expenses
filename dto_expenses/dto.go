package dto

import "time"

type ExpensesType struct {
	Id    int64
	Title string
}

type Expense struct {
	Id            int       `json:"id"`
	ExpenseTypeId string    `json:"expense_type_id" db:"expense_type_id"`
	Time          time.Time `json:"time" db:"reated_at"`
	SpentMoney    float64   `json:"spent_money" db:"spent_money"`
}

type CreateExpense struct {
	ExpenseType string  `json:"expense_type" binding:"required"`
	Time        string  `json:"time" binding:"required"`
	SpentMoney  float64 `json:"spent_money" binding:"required"`
}
type UpdateExpense struct {
	Time       string  `json:"time"`
	SpentMoney float64 `json:"spent_money"`
}

type TypesExpenseUserParams struct {
	Name  string `json:"name"`
	Login string `json:"login"`
	Id    int    `json:"id"`
}

type User struct {
	Id      int    `json:"-"`
	Name    string `json:"name" binding:"required"`
	Surname string `json:"surname" binding:"required"`
	Login   string `json:"login" binding:"required"`
	Pass    string `json:"pass" binding:"required"`
	Email   string `json:"email" binding:"required"`
}
