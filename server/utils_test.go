package server

import (
	"context"
	"expenses/app/auth"
	dto "expenses/dto_expenses"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) GetTypesExpenseUser(ctx context.Context, userId int) ([]dto.ExpensesType, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).([]dto.ExpensesType), args.Error(1)
}

func (m *MockRepo) GetUserId(ctx context.Context, expenseID int) (int, error) {
	args := m.Called(ctx, expenseID)
	return args.Int(0), args.Error(1)
}

func (m *MockRepo) IsExpenseTypeExists(ctx context.Context, expType string) (bool, error) {
	args := m.Called(ctx, expType)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepo) IsExpenseExists(ctx context.Context, expenseID int) (bool, error) {
	args := m.Called(ctx, expenseID)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepo) CreateExpenseType(ctx context.Context, tx pgx.Tx, expType string, userId int) (int, error) {
	args := m.Called(ctx, tx, expType, userId)
	return args.Int(0), args.Error(1)
}

func (m *MockRepo) GetExpenseTypeID(ctx context.Context, tx pgx.Tx, expType string) (int, error) {
	args := m.Called(ctx, tx, expType)
	return args.Int(0), args.Error(1)
}

func (m *MockRepo) SetExpenseTimeAndSpent(ctx context.Context, tx pgx.Tx, expTypeId int, timeSpent string, spent float64) (int, error) {
	args := m.Called(ctx, tx, expTypeId, timeSpent, spent)
	return args.Int(0), args.Error(1)
}

func (m *MockRepo) AddFileExpense(ctx context.Context, filepath string, expId int, typeFile string) error {
	args := m.Called(ctx, filepath, expId, typeFile)
	return args.Error(0)
}

func (m *MockRepo) CreateUserExpense(ctx context.Context, expenseData *dto.CreateExpense, userId int) (int, error) {
	args := m.Called(ctx, expenseData, userId)
	return args.Int(0), args.Error(1)
}

func (m *MockRepo) GetAllExpenses(ctx context.Context, userId int) ([]dto.Expense, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).([]dto.Expense), args.Error(1)
}

func (m *MockRepo) DeleteExpense(ctx context.Context, expenseId, userId int) (int, error) {
	args := m.Called(ctx, expenseId, userId)
	return args.Int(0), args.Error(1)
}

func (m *MockRepo) DeleteFile(ctx context.Context, pathFile string, expenseId int) error {
	args := m.Called(ctx, pathFile, expenseId)
	return args.Error(0)
}

func (m *MockRepo) GetExpense(ctx context.Context, userID, expenseID int) (*dto.Expense, error) {
	args := m.Called(ctx, userID, expenseID)
	return args.Get(0).(*dto.Expense), args.Error(1)
}

func (m *MockRepo) UpdateExpense(ctx context.Context, expenseID int, newExpense *dto.Expense) error {
	args := m.Called(ctx, expenseID, newExpense)
	return args.Error(0)
}

func (m *MockRepo) CreateUser(ctx context.Context, user *dto.User) (int, error) {
	args := m.Called(ctx, user)
	return args.Int(0), args.Error(1)
}

func (m *MockRepo) GetUser(userName, hashPassword string) (*dto.User, error) {
	args := m.Called(userName, hashPassword)
	return args.Get(0).(*dto.User), args.Error(1)
}

func TestGenerateToken(t *testing.T) {

	mockRepository := new(MockRepo)
	handler := &Handler{repo: mockRepository}

	name := "AidAmir"
	pass := "qwerty"
	hashPass := auth.HashPassword(pass)
	user := &dto.User{Id: 111}

	mockRepository.On("GetUser", name, hashPass).Return(user, nil)

	tokenString, err := auth.GenerateToken(handler.repo, name, pass)

	assert.NotEmpty(t, tokenString)
	assert.NoError(t, err)

	mockRepository.AssertCalled(t, "GetUser", name, hashPass)

	token, err := jwt.ParseWithClaims(tokenString, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(signingKey), nil
	})
	assert.NoError(t, err)

	claims, ok := token.Claims.(*tokenClaims)
	assert.True(t, ok)
	assert.Equal(t, user.Id, claims.UserId)

	mockRepository.AssertExpectations(t)

}

func TestParseToken(t *testing.T) {
	expectedID := 123
	tokenSigned, err := generateTestToken(t, expectedID)
	assert.NoError(t, err)

	userID, err := auth.ParseToken(tokenSigned)

	assert.NoError(t, err)
	assert.Equal(t, expectedID, userID)
}

func generateTestToken(t testing.TB, userID int) (string, error) {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		userID,
	})
	return token.SignedString([]byte(signingKey))
}
