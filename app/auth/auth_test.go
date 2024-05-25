package auth

import (
	dto "expenses/dto_expenses"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetUser(userName, hashPassword string) (*dto.User, error) {
	args := m.Called(userName, hashPassword)
	return args.Get(0).(*dto.User), args.Error(1)
}

func TestParseToken(t *testing.T) {
	expectedID := 123
	tokenSigned, err := generateTestToken(t, expectedID)
	assert.NoError(t, err)

	userID, err := parseToken(tokenSigned)

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

func TestGenerateToken(t *testing.T) {

	mockRepository := new(MockRepository)

	name := "AidAmir"
	pass := "qwerty"
	hashPass := hashPassword(pass)
	user := &dto.User{Id: 111}

	mockRepository.On("GetUser", name, hashPass).Return(user, nil)

	tokenString, err := generateToken(mockRepository, name, pass)

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
