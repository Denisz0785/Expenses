package dto

type Expenses struct {
	Id    int64
	Title string
}

type TypesExpenseUserParams struct {
	Name  string
	Login string
	Id    int
}
