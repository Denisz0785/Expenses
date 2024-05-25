// Package repository contains methods to work with database
package repository

import (
	"context"
	"errors"
	dto "expenses/dto_expenses"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

type Repository interface {
	GetTypesExpenseUser(ctx context.Context, userId int) ([]dto.ExpensesType, error)
	GetUserId(ctx context.Context, expenseID int) (int, error)
	IsExpenseTypeExists(ctx context.Context, expType string) (bool, error)
	IsExpenseExists(ctx context.Context, expenseID int) (bool, error)
	CreateExpenseType(ctx context.Context, tx pgx.Tx, expType string, userId int) (int, error)
	GetExpenseTypeID(ctx context.Context, tx pgx.Tx, expType string) (int, error)
	SetExpenseTimeAndSpent(ctx context.Context, tx pgx.Tx, expTypeId int, timeSpent string, spent float64) (int, error)
	AddFileExpense(ctx context.Context, filepath string, expId int, typeFile string) error
	CreateUserExpense(ctx context.Context, expenseData *dto.CreateExpense, userId int) (int, error)
	GetAllExpenses(ctx context.Context, userId int) ([]dto.Expense, error)
	DeleteExpense(ctx context.Context, expenseId, userId int) (int, error)
	DeleteFile(ctx context.Context, pathFile string, expenseId int) error
	GetExpense(ctx context.Context, userID, expenseID int) (*dto.Expense, error)
	UpdateExpense(ctx context.Context, expenseID int, newExpense *dto.Expense) error
	CreateUser(ctx context.Context, user *dto.User) (int, error)
	GetUser(userName, hashPassword string) (*dto.User, error)
}

// ExpenseRepo create custom struct which contains descriptor of connection to database
type ExpenseRepo struct {
	conn *pgx.Conn
}

// NewExpenseRepo create ExpenseRepo
func NewExpenseRepo(conn *pgx.Conn) *ExpenseRepo {
	return &ExpenseRepo{conn: conn}
}

// GetTypesExpenseUser get types of expenses from database by users's od or name or login
func (r *ExpenseRepo) GetTypesExpenseUser(ctx context.Context, userId int) ([]dto.ExpensesType, error) {
	var query string
	query = fmt.Sprint("SELECT title, id from expense_type where users_id=$1")
	/*
		if d.Id == 0 && d.Login == "" && d.Name == "" {
			err := errors.New("incorrect user data")
			return nil, err
		}

		sql := "SELECT e.title, e.id from expense_type e"
		sqlWhere := ",users where e.users_id=users.id and users."
		var param interface{}
		if d.Id != 0 {
			query = fmt.Sprint(sql + " where e.users_id=$1")
			param = d.Id
		} else if d.Login != "" {
			query = fmt.Sprint(sql + sqlWhere + "login=$1")
			param = d.Login
		} else if d.Name != "" {
			query = fmt.Sprint(sql + sqlWhere + "name=$1")
			param = d.Name
		}
	*/

	rows, _ := r.conn.Query(ctx, query, userId)
	expense, err := pgx.CollectRows(rows, pgx.RowToStructByName[dto.ExpensesType])
	if err != nil {
		err = fmt.Errorf("unable to connect to database: %v", err)
		return nil, err
	}
	return expense, nil
}

// GetUserId get id of user from database by id of expense
func (r *ExpenseRepo) GetUserId(ctx context.Context, expenseID int) (int, error) {
	var userId int
	query := fmt.Sprintf("SELECT et.users_id from expense_type et join expense e on et.id=e.expense_type_id where e.id=%v;", expenseID)
	err := r.conn.QueryRow(ctx, query).Scan(&userId)
	if err != nil {
		fmt.Println(err.Error())
		return 0, err
	}
	return userId, nil
}

// IsExpenseTypeExists checking exist type of expense or not in a database
func (r *ExpenseRepo) IsExpenseTypeExists(ctx context.Context, expType string) (bool, error) {
	existExpense := false

	err := r.conn.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM expense_type WHERE title=$1);", expType).Scan(&existExpense)
	if err != nil {
		log.Println(err)
		return existExpense, err
	}
	return existExpense, nil
	/*
		rows, _ := r.conn.Query(ctx, "Select title from expense_type")
		numbers, err := pgx.CollectRows(rows, pgx.RowTo[string])
		existExpType := false
		if err != nil {
			return existExpType, err
		} else {
			for _, v := range numbers {
				if v == *expType {
					existExpType = true
					return existExpType, nil
				}

			}
			existExpType = false
		}
		return existExpType, nil

	*/
}

// IsExpenseExists checks existing expense in a database by expense's id
func (r *ExpenseRepo) IsExpenseExists(ctx context.Context, expenseID int) (bool, error) {
	existExpense := false

	err := r.conn.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM expense WHERE id=$1);", expenseID).Scan(&existExpense)
	if err != nil {
		fmt.Println(err.Error())
		return existExpense, err
	}
	return existExpense, nil
}

// CreateExpenseType insert a new type of expenses in a table expense_type
func (r *ExpenseRepo) CreateExpenseType(ctx context.Context, tx pgx.Tx, expType string, userId int) (int, error) {
	var expenseTypeId int
	err := tx.QueryRow(ctx, "Insert into expense_type(users_id,title) values ($1,$2) returning id", userId, expType).Scan(&expenseTypeId)
	if err != nil {
		log.Println(err)
		return -1, err
	}
	return expenseTypeId, nil
}

// GetExpenseTypeID gets id of expense type
func (r *ExpenseRepo) GetExpenseTypeID(ctx context.Context, tx pgx.Tx, expType string) (int, error) {
	var expTypeId int
	err := tx.QueryRow(ctx, "select id from expense_type where title=$1", expType).Scan(&expTypeId)
	if err != nil {
		log.Println(err)
		return -1, err
	}
	return expTypeId, err
}

// SetExpenseTimeAndSpent Creates new row in a expense table
func (r *ExpenseRepo) SetExpenseTimeAndSpent(ctx context.Context, tx pgx.Tx, expTypeId int, timeSpent string, spent float64) (int, error) {
	// Create a new row into table expense
	var expenseId int
	err := tx.QueryRow(ctx, "Insert into expense(expense_type_id,reated_at, spent_money) values ($1,$2,$3) returning id", expTypeId, timeSpent, spent).Scan(&expenseId)
	if err != nil {
		log.Println(err)
		return -1, err
	}
	return expenseId, nil
}

// AddFileExpense define type of the file and write info of file to the database
func (r *ExpenseRepo) AddFileExpense(ctx context.Context, filepath string, expId int, typeFile string) error {

	query := "INSERT INTO files (expense_id,path_file, type_file) VALUES ($1,$2,$3)"
	_, err := r.conn.Exec(ctx, query, expId, filepath, typeFile)
	if err != nil {
		return err
	}

	return nil
}

// CreateUserExpense checks existing type of expenses from command-line in a table, and Creates new row to expense table by transaction
func (r *ExpenseRepo) CreateUserExpense(ctx context.Context, expenseData *dto.CreateExpense, userId int) (int, error) {
	var existExpType bool
	var expenseId int
	var err error
	// checking expType exists in a table expense_type or not
	existExpType, err = r.IsExpenseTypeExists(ctx, expenseData.ExpenseType)
	if err != nil {
		log.Println(err)
		return -1, err
	}

	// begin transaction
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return -1, err
	}
	defer tx.Commit(ctx)

	if !existExpType {
		/*
			var userId int
			loginValue := *login
			// by QueryRow gets user's id from table users by login
			err = r.conn.QueryRow(ctx,
				`SELECT id FROM users where login=$1`, loginValue).Scan(&userId)
			if err != nil {
				err = fmt.Errorf("QueryRow failed: %v", err)
				return err
			}

		*/
		// Create new type of expense to expense_type table
		expId, err := r.CreateExpenseType(ctx, tx, expenseData.ExpenseType, userId)
		if err != nil {
			log.Println(err)
			return -1, err
		}
		/*
			// getting Id of new expense_type
			expId, err1 := r.GetExpenseTypeID(ctx, tx, expType)
			if err1 != nil {
				return -1,err1
			}

		*/
		// Create new expense
		expenseId, err = r.SetExpenseTimeAndSpent(ctx, tx, expId, expenseData.Time, expenseData.SpentMoney)
		if err != nil {
			return -1, err
		}

	} else if existExpType {
		// getting Id of expType from expense_type
		expId, err := r.GetExpenseTypeID(ctx, tx, expenseData.ExpenseType)
		if err != nil {
			log.Println(err)
			return -1, err
		}
		// Creating a new row into expense table
		expenseId, err = r.SetExpenseTimeAndSpent(ctx, tx, expId, expenseData.Time, expenseData.SpentMoney)
		if err != nil {
			log.Println(err)
			return -1, err
		}

	}
	return expenseId, nil
}

func (r *ExpenseRepo) GetAllExpenses(ctx context.Context, userId int) ([]dto.Expense, error) {

	query := fmt.Sprint(`select e.id,e.expense_type_id,e.reated_at,e.spent_money from users u
		join  expense_type et on u.id=et.users_id join expense e on e.expense_type_id=et.id where u.id=$1`)
	rows, err := r.conn.Query(ctx, query, userId)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	expenses, err := pgx.CollectRows(rows, pgx.RowToStructByName[dto.Expense])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return expenses, nil
}

func (r *ExpenseRepo) DeleteExpense(ctx context.Context, expenseId, userId int) (int, error) {
	var idDeleteExpense int
	query := `DELETE FROM expense WHERE id IN (select e.id from users u join  expense_type et
		ON u.id=et.users_id join expense e on e.expense_type_id=et.id where u.id=$1
		and e.id=$2) returning id`
	//_, err := r.conn.Exec(ctx, query, userId, expenseId)
	err := r.conn.QueryRow(ctx, query, userId, expenseId).Scan(&idDeleteExpense)
	if err != nil {
		if idDeleteExpense == 0 {
			log.Println("id expense does not exist")
			return -1, errors.New("expense does not exist")
		}
		log.Println(err)
		return -1, err
	}
	return idDeleteExpense, nil
}

// ConnectToDB connects to DB
func ConnectToDB(ctx context.Context, myurl string) (*pgx.Conn, error) {
	fmt.Println(os.Getenv(myurl))
	conn, err := pgx.Connect(ctx, os.Getenv(myurl))
	if err != nil {
		err = fmt.Errorf("unable to connect to database: %v", err)
		return nil, err
	}
	return conn, nil
}

// DeleteFile removes file from database by name of the file and expense Id
func (r *ExpenseRepo) DeleteFile(ctx context.Context, pathFile string, expenseId int) error {
	query := fmt.Sprintf("DELETE FROM files WHERE path_file='%v' AND expense_id=%v;", pathFile, expenseId)

	_, err := r.conn.Exec(ctx, query)
	if err != nil {
		return err
	}
	return nil
}

func (r *ExpenseRepo) GetExpense(ctx context.Context, userID, expenseID int) (*dto.Expense, error) {
	var expense dto.Expense
	query := `select e.id,e.expense_type_id,e.reated_at,e.spent_money from users u
		join  expense_type et on u.id=et.users_id join expense e on e.expense_type_id=et.id where u.id=$1 and e.id=$2`
	err := r.conn.QueryRow(ctx, query, userID, expenseID).Scan(&expense.Id, &expense.ExpenseTypeId,
		&expense.Time, &expense.SpentMoney)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &expense, nil
}

func (r *ExpenseRepo) UpdateExpense(ctx context.Context, expenseID int, newExpense *dto.Expense) error {
	query := "Update expense set reated_at=$1, spent_money=$2 where id=$3"
	_, err := r.conn.Exec(ctx, query, newExpense.Time, newExpense.SpentMoney, expenseID)
	if err != nil {
		return err
	}
	return nil
}
