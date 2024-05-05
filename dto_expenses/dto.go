package dto

type Expenses struct {
	Id    int64
	Title string
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
