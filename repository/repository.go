package repository

import (
	"context"
	dto "expenses/dto_expenses"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

type repo interface {
	GetExpenseType()
	CreateValuesDB()
	GetUserExpenseTypes()
	CreateExpense()
	CreateExpenseType()
	GetExpenseTypeID()
	SetExpenseTimeAndSpent()
}

type ExpenseRepo struct {
	conn *pgx.Conn
}

func NewExpenseRepo(conn *pgx.Conn) *ExpenseRepo {
	return &ExpenseRepo{conn: conn}
}

func (r *ExpenseRepo) GetTypesExpenseUser(ctx context.Context, id1 int) ([]dto.Expenses, error) {
	rows, _ := r.conn.Query(ctx, "SELECT id, title from expense_type where expense_type.users_id=$1", id1)
	expense, err := pgx.CollectRows(rows, pgx.RowToStructByName[dto.Expenses])
	if err != nil {
		err = fmt.Errorf("unable to connect to database: %v", err)
		return nil, err
	}
	return expense, nil
}

// GetUserExpenseTypes gets all rows of type of expenses from DB by name
func (r *ExpenseRepo) GetUserExpenseTypes(ctx context.Context, name string) ([]string, error) {
	rows, _ := r.conn.Query(ctx, "SELECT title from expense_type, users where expense_type.users_id=users.id and users.name=$1", name)
	numbers, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		err = fmt.Errorf("unable to connect to database: %v", err)
		return nil, err
	}
	return numbers, nil
}

// GetManyRows gets all rows of type of expenses from DB by login
func (r *ExpenseRepo) GetExpenseTypesUser(ctx context.Context, login string) ([]string, error) {
	rows, _ := r.conn.Query(ctx, "SELECT title from expense_type, users where expense_type.users_id=users.id and users.login=$1", login)
	numbers, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		err = fmt.Errorf("unable to connect to database: %v", err)
		return nil, err
	}
	return numbers, nil
}

// IsExpenseTypeExists checking exist type of expense or not in a database
func (r *ExpenseRepo) IsExpenseTypeExists(ctx context.Context, expType *string) (bool, error) {
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
}

// CreateExpenseType insert a new type of expenses in a table expense_type
func (r *ExpenseRepo) CreateExpenseType(ctx context.Context, tx pgx.Tx, expType *string, userId int) error {
	_, err := tx.Exec(ctx, "Insert into expense_type(users_id,title) values ($1,$2)", userId, *expType)
	if err != nil {
		return err
	}
	return err
}

// GetExpenseTypeID gets id of expense type
func (r *ExpenseRepo) GetExpenseTypeID(ctx context.Context, tx pgx.Tx, expType *string) (*int, error) {
	var expTypeId int
	err := tx.QueryRow(ctx, "select id from expense_type where title=$1", *expType).Scan(&expTypeId)
	if err != nil {
		err = fmt.Errorf("QueryRow failed: %v", err)
		return nil, err
	}
	return &expTypeId, err
}

// SetExpenseTimeAndSpent Creates new row in a expense table
func (r *ExpenseRepo) SetExpenseTimeAndSpent(ctx context.Context, tx pgx.Tx, expTypeId *int, timeSpent *string, spent *float64) error {
	// Create a new row into table expense

	_, err := tx.Exec(ctx, "Insert into expense(expense_type_id,reated_at, spent_money) values ($1,$2,$3)", *expTypeId, *timeSpent, *spent)
	if err != nil {
		return err
	}
	fmt.Println("was Createed")
	return err
}

// CreateUserExpense checks existing type of expenses from command-line in a table, and Creates new row to expense table by transaction
func (r *ExpenseRepo) CreateUserExpense(ctx context.Context, login *string, expType *string, timeSpent *string, spent *float64) error {
	// checking expType exists in a table expense_type or not
	existExpType, err := r.IsExpenseTypeExists(ctx, expType)
	if err != nil {
		return err
	}

	// begin transaction
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if !existExpType {
		var userId int
		loginValue := *login
		// by QueryRow gets user's id from table users by login
		err = r.conn.QueryRow(ctx,
			`SELECT id FROM users where login=$1`, loginValue).Scan(&userId)
		if err != nil {
			err = fmt.Errorf("QueryRow failed: %v", err)
			return err
		}
		// Createing new type of expense to expense_type table
		err = r.CreateExpenseType(ctx, tx, expType, userId)
		if err != nil {
			return err
		}
		// getting Id of new expense_type
		expId, err1 := r.GetExpenseTypeID(ctx, tx, expType)
		if err1 != nil {
			return err1
		}
		// Createing a new row into expense table
		err = r.SetExpenseTimeAndSpent(ctx, tx, expId, timeSpent, spent)
		if err != nil {
			return err
		}

	} else if existExpType {
		// getting Id of expType from expense_type
		expId, err1 := r.GetExpenseTypeID(ctx, tx, expType)
		if err1 != nil {
			return err1
		}
		// Createing a new row into expense table
		err = r.SetExpenseTimeAndSpent(ctx, tx, expId, timeSpent, spent)
		if err != nil {
			return err
		}

	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return err
}

// ConnectToDB connects to DB
func ConnectToDB(ctx context.Context, myurl string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, os.Getenv(myurl))
	if err != nil {
		err = fmt.Errorf("unable to connect to database: %v", err)
		return nil, err
	}
	return conn, nil
}
