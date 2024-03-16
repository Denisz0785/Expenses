package repository

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/jackc/pgx/v5"
)

// GetExpenseType gets one row of type of expenses from DB by name
func GetExpenseType(conn *pgx.Conn, name string) (*string, error) {

	var typeExpenses string
	err := conn.QueryRow(context.Background(),
		"SELECT type_expenses from expense_type, users where expense_type.users_id=users.id and users.name=$1", name).Scan(&typeExpenses)
	if err != nil {
		err1 := fmt.Errorf("unable to connect to database: %v", err)
		return nil, err1
	}
	return &typeExpenses, nil
}

// GetManyRows gets all rows of type of expenses from DB by name
func GetManyRowsByName(conn *pgx.Conn, name string) ([]string, error) {
	rows, _ := conn.Query(context.Background(), "SELECT type_expenses from expense_type, users where expense_type.users_id=users.id and users.name=$1", name)
	numbers, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		err1 := fmt.Errorf("unable to connect to database: %v", err)
		return nil, err1
	}
	return numbers, nil
}

// GetManyRows gets all rows of type of expenses from DB by login
func GetManyRowsByLogin(conn *pgx.Conn, login string) ([]string, error) {
	rows, _ := conn.Query(context.Background(), "SELECT type_expenses from expense_type, users where expense_type.users_id=users.id and users.login=$1", login)
	numbers, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		err1 := fmt.Errorf("unable to connect to database: %v", err)
		return nil, err1
	}
	return numbers, nil
}

// AddValuesDB insert row to the table user
func AddValuesDB(conn *pgx.Conn) error {

	comTag, err := conn.Exec(context.Background(), "Insert into users(name,surname,login,pass,email) VALUES ('kolya','Bon','spiman','1243','er1@23.ru')")
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
		err1 := fmt.Errorf("unable to connect to database: %v", err)
		return nil, err1
	}
	return conn, nil
}

// AddExpense checking type of expenses in a database, if not exists - adds new type of expenses to a table expense_type
// after this adds a new row to a table expense
func AddExpense(conn *pgx.Conn, login *string, expType *string, timeSpent *string, spent *float64) error {
	// getting all typies of expenses by login and checking expType there or not
	numbers, err := GetManyRowsByLogin(conn, *login)
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
			err = conn.QueryRow(context.Background(), "SELECT id FROM users where login=$1", loginValue).Scan(&userId)
			fmt.Println(userId)
			if err != nil {
				err1 := fmt.Errorf("QueryRow failed: %v", err)
				return err1
			}
			// bigin transaction
			tx, err := conn.Begin(context.Background())
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
			err = conn.QueryRow(context.Background(), "select id from expense_type where type_expenses=$1", *expType).Scan(&expTypeId)
			if err != nil {
				err1 := fmt.Errorf("QueryRow failed: %v", err)
				return err1
			}
			// delete all "-" from timeSpent
			FormatTimeSpent := strings.Replace(*timeSpent, "-", " ", 3)
			// add a new row into table expense
			_, err = tx.Exec(context.Background(), "Insert into expense(expense_type_id,reated_at, spent_money) values ($1,$2,$3)", expTypeId, FormatTimeSpent, *spent)
			if err != nil {
				return err
			}

			err = tx.Commit(context.Background())
			if err != nil {
				return err
			}
			// if expType exist in a table expense_type
		} else if existExpType {
			var expTypeId int
			// by QueryRow gets id of expType from type of expense
			err = conn.QueryRow(context.Background(), "select id from expense_type where type_expenses=$1", *expType).Scan(&expTypeId)
			if err != nil {
				err1 := fmt.Errorf("QueryRow failed: %v", err)
				return err1
			}
			// delete all "-" from timeSpent
			FormatTimeSpent := strings.Replace(*timeSpent, "-", " ", 3)
			// add a new row into table expense
			com, err := conn.Exec(context.Background(), "Insert into expense(expense_type_id,reated_at, spent_money) values ($1,$2,$3)", expTypeId, FormatTimeSpent, *spent)
			if com.RowsAffected() != 1 {
				return errors.New("add is not done")
			}
			return err
		}

	}
	return nil
}
