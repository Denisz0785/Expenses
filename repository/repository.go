package repository

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

type repository interface {
	GetExpenseType()
	AddValuesDB()
	GetManyRowsByName()
	AddExpense()
	AddExpType()
	GetIdTypeExp()
	AddExpValue()
}

type ExpenseRepo struct {
	conn *pgx.Conn
	tx1  pgx.Tx
}

func NewExpenseRepo(conn *pgx.Conn) *ExpenseRepo {
	return &ExpenseRepo{conn: conn}
}

// GetExpenseType gets one row of type of expenses from DB by name
func (r *ExpenseRepo) GetExpenseType(name string) (*string, error) {

	var typeExpenses string
	err := r.conn.QueryRow(context.Background(),
		"SELECT type_expenses from expense_type, users where expense_type.users_id=users.id and users.name=$1", name).Scan(&typeExpenses)
	if err != nil {
		err = fmt.Errorf("unable to connect to database: %v", err)
		return nil, err
	}
	return &typeExpenses, nil
}

// GetManyRows gets all rows of type of expenses from DB by name
func (r *ExpenseRepo) GetManyRowsByName(name string) ([]string, error) {
	rows, _ := r.conn.Query(context.Background(), "SELECT type_expenses from expense_type, users where expense_type.users_id=users.id and users.name=$1", name)
	numbers, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		err = fmt.Errorf("unable to connect to database: %v", err)
		return nil, err
	}
	return numbers, nil
}

// GetManyRows gets all rows of type of expenses from DB by login
func (r *ExpenseRepo) GetManyRowsByLogin(login string) ([]string, error) {
	rows, _ := r.conn.Query(context.Background(), "SELECT type_expenses from expense_type, users where expense_type.users_id=users.id and users.login=$1", login)
	numbers, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		err = fmt.Errorf("unable to connect to database: %v", err)
		return nil, err
	}
	return numbers, nil
}

// CheckExistTypeExp checking exist type of expense or not in a database
func (r *ExpenseRepo) CheckExistTypeExp(expType *string) (bool, error) {
	rows, _ := r.conn.Query(context.Background(), "Select type_expenses from expense_type")
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

// AddExpType insert a new type of expenses in a table expense_type
func (r *ExpenseRepo) AddExpType(expType *string, userId int) error {
	_, err := r.tx1.Exec(context.Background(), "Insert into expense_type(users_id,type_expenses) values ($1,$2)", userId, *expType)
	if err != nil {
		return err
	}
	return err
}

// GetIdTypeExp gets id of expense type
func (r *ExpenseRepo) GetIdTypeExp(expType *string) (*int, error) {
	var expTypeId int
	err := r.tx1.QueryRow(context.Background(), "select id from expense_type where type_expenses=$1", *expType).Scan(&expTypeId)
	if err != nil {
		err = fmt.Errorf("QueryRow failed: %v", err)
		return nil, err
	}
	return &expTypeId, err
}

// AddExpValue adds new row in a expense table
func (r *ExpenseRepo) AddExpValue(expTypeId *int, timeSpent *string, spent *float64) error {
	// add a new row into table expense

	_, err := r.tx1.Exec(context.Background(), "Insert into expense(expense_type_id,reated_at, spent_money) values ($1,$2,$3)", *expTypeId, *timeSpent, *spent)
	if err != nil {
		return err
	}
	fmt.Println("was added")
	return err
}

// AddExpTransaction checks existing type of expenses from command-line in a table, and adds new row to expense table by transaction
func (r *ExpenseRepo) AddExpTransaction(login *string, expType *string, timeSpent *string, spent *float64) error {
	// checking expType exists in a table expense_type or not
	existExpType, err := r.CheckExistTypeExp(expType)
	if err != nil {
		return err
	}

	// begin transaction
	tx, err := r.conn.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	r.tx1 = tx

	if !existExpType {
		var userId int
		loginValue := *login
		// by QueryRow gets user's id from table users by login
		err = r.conn.QueryRow(context.Background(),
			`SELECT id FROM users where login=$1`, loginValue).Scan(&userId)
		if err != nil {
			err = fmt.Errorf("QueryRow failed: %v", err)
			return err
		}
		// adding new type of expense to expense_type table
		err = r.AddExpType(expType, userId)
		if err != nil {
			return err
		}
		// getting Id of new expense_type
		expId, err1 := r.GetIdTypeExp(expType)
		if err1 != nil {
			return err1
		}
		// adding a new row into expense table
		err = r.AddExpValue(expId, timeSpent, spent)
		if err != nil {
			return err
		}

		fmt.Println("new expense was added")

	} else if existExpType {
		// getting Id of expType from expense_type
		expId, err1 := r.GetIdTypeExp(expType)
		if err1 != nil {
			return err1
		}
		// adding a new row into expense table
		err = r.AddExpValue(expId, timeSpent, spent)
		if err != nil {
			return err
		}

	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}
	return err
}

// AddValuesDB insert row to the table user
func (r *ExpenseRepo) AddValuesDB() error {

	comTag, err := r.conn.Exec(context.Background(), "Insert into users(name,surname,login,pass,email) VALUES ('kolya','Bon','spiman','1243','er1@23.ru')")
	if err != nil {
		return err
	}
	if comTag.RowsAffected() != 1 {
		return errors.New("add is not done")
	}
	return nil
}

// ConnectToDB connects to DB
func ConnectToDB(myurl string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), os.Getenv(myurl))
	if err != nil {
		err = fmt.Errorf("unable to connect to database: %v", err)
		return nil, err
	}
	return conn, nil
}

// AddExpense checking type of expenses in a database, if not exists - adds new type of expenses to a table expense_type
// after this adds a new row to a table expense
func (r *ExpenseRepo) AddExpense(login *string, expType *string, timeSpent *string, spent *float64) error {
	// getting all typies of expenses by login and checking expType there or not
	numbers, err := r.GetManyRowsByLogin(*login)
	if err != nil {
		return err
	} else {
		existExpType := false
		for _, v := range numbers {
			if v == *expType {
				existExpType = true
				break
			}
		}
		// if expType not exists in a table expense_type
		if !existExpType {
			var userId int
			loginValue := *login
			// by QueryRow gets user's id from table users by login
			err = r.conn.QueryRow(context.Background(),
				`SELECT id 
			FROM 
			users 
			where login=$1`, loginValue).Scan(&userId)
			if err != nil {
				err = fmt.Errorf("QueryRow failed: %v", err)
				return err
			}
			// begin transaction
			tx, err := r.conn.Begin(context.Background())
			if err != nil {
				return err
			}
			defer tx.Rollback(context.Background())
			// insert a new type of expenses in a table expense_type
			_, err = tx.Exec(context.Background(), "Insert into expense_type(users_id,type_expenses) values ($1,$2)", userId, *expType)
			if err != nil {
				return err
			}
			// by QueryRow gets id expType from a table expense_type
			var expTypeId int
			err = r.conn.QueryRow(context.Background(), "select id from expense_type where type_expenses=$1", *expType).Scan(&expTypeId)
			if err != nil {
				err = fmt.Errorf("QueryRow failed: %v", err)
				return err
			}
			// delete all "-" from timeSpent
			// FormatTimeSpent := strings.Replace(*timeSpent, "-", " ", 3)
			// add a new row into table expense
			_, err = tx.Exec(context.Background(), "Insert into expense(expense_type_id,reated_at, spent_money) values ($1,$2,$3)", expTypeId, *timeSpent, *spent)
			if err != nil {
				return err
			}

			err = tx.Commit(context.Background())
			if err != nil {
				return err
			}
			fmt.Println("new expense was added")
			// if expType exist in a table expense_type
		} else if existExpType {
			var expTypeId int
			// by QueryRow gets id of expType from type of expense
			err = r.conn.QueryRow(context.Background(), "select id from expense_type where type_expenses=$1", *expType).Scan(&expTypeId)
			if err != nil {
				err = fmt.Errorf("QueryRow failed: %v", err)
				return err
			}
			// delete all "-" from timeSpent
			// FormatTimeSpent := strings.Replace(*timeSpent, "-", " ", 4)
			// add a new row into table expense
			com, err := r.conn.Exec(context.Background(), "Insert into expense(expense_type_id,reated_at, spent_money) values ($1,$2,$3)", expTypeId, *timeSpent, *spent)
			if com.RowsAffected() != 1 {
				return errors.New("add is not done")
			} else {
				fmt.Println("new expense was added")
			}
			return err
		}

	}
	return nil
}
