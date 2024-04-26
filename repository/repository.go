package repository

import (
	"context"
	"errors"
	dto "expenses/dto_expenses"
	"fmt"
	"os"
	"strings"

	"github.com/jackc/pgx/v5"
)

type Repository interface {
	GetTypesExpenseUser()
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

func (r *ExpenseRepo) GetTypesExpenseUser(ctx context.Context, d *dto.TypesExpenseUserParams) ([]dto.Expenses, error) {
	if d.Id == 0 && d.Login == "" && d.Name == "" {
		err := errors.New("can't find info abour User")
		return nil, err
	}

	var query string

	if d.Id != 0 {
		query = fmt.Sprintf("SELECT id, title from expense_type where expense_type.users_id=%d", d.Id)

	} else if d.Login != "" {
		query = fmt.Sprintf("SELECT e.title, e.id from expense_type e, users where e.users_id=users.id and users.login=%s", d.Login)

	} else if d.Name != "" {
		query = fmt.Sprintf("SELECT e.title, e.id from expense_type e, users where e.users_id=users.id and users.name=%s", d.Name)
	}
	rows, _ := r.conn.Query(ctx, query)
	expense, err := pgx.CollectRows(rows, pgx.RowToStructByName[dto.Expenses])
	if err != nil {
		err = fmt.Errorf("unable to connect to database: %v", err)
		return nil, err
	}
	return expense, nil
}

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
	fmt.Println("was Created")
	return err
}

func (r *ExpenseRepo) AddFileExpense(ctx context.Context, filepath string, expId int) error {
	var typeFile string
	src := strings.Split(filepath, ".")
	if len(src) == 2 {
		extension := strings.ToLower(src[1])
		switch extension {
		case "doc", "pdf", "txt":
			typeFile = "document"
		case "jpg", "jpeg", "png", "gif", "raw", "svg", "bmp", "ico", "tiff", "webp":
			typeFile = "image"
		case "mp4", "webm", "mov", "avi", "flv", "wmv", "mkv", "mpeg", "3gp", "ogv":
			typeFile = "video"
		default:
			return errors.New("формат сохраняемого файла не поддерживается")
		}
	}
	query := "INSERT INTO files (expense_id,path_file, type_file) VALUES ($1,$2,$3)"
	_, err := r.conn.Exec(ctx, query, expId, filepath, typeFile)
	if err != nil {
		return err
	}
	return nil
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

func (r *ExpenseRepo) DeleteFile(ctx context.Context, pathFile string, expenseId int) error {
	query := fmt.Sprintf("DELETE FROM files WHERE path_file='%v' AND expense_id=%v;", pathFile, expenseId)

	_, err := r.conn.Exec(ctx, query)
	if err != nil {
		return err
	}
	return nil
}
