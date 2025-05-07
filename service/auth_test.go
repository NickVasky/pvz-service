package service

import (
	"AvitoTechPVZ/repo"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGenerateJwtToken(t *testing.T) {
	// Mock user and role
	user := repo.User{
		Id: uuid.New(),
	}
	role := repo.Role{
		Name: "admin",
	}

	claims := NewUserClaims(user, role)
	// Generate token
	token := claims.GenerateJwtToken()

	// Assert no error
	assert.NotEmpty(t, token)

	parsedClaims := &UserClaims{}
	// Parse and validate token
	parsedToken, err := jwt.ParseWithClaims(token, parsedClaims, jwtKeyFunc)

	// Assert no error during parsing
	assert.NoError(t, err)
	assert.NotNil(t, parsedToken)

	// Assert claims
	if parsedToken.Valid {
		parsedUuid, _ := uuid.Parse(parsedClaims.UserID)
		assert.Equal(t, user.Id, parsedUuid)
		assert.Equal(t, role.Name, parsedClaims.Role)
	} else {
		t.Fatalf("invalid token claims")
	}
}

func TestParseUserToken(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNzhiZTg4NzEtOTBlYi00YjVmLWI3ZDYtMjI4ZmYyOGYyM2MwIiwicm9sZSI6ImFkbWluIn0.P-pnVm4lZ01QAI_cS2wYNVhs8XZJUp0a3lHHsc78cAA"
	expectedUserId := "78be8871-90eb-4b5f-b7d6-228ff28f23c0"
	expectedRole := "admin"
	claims, err := ParseUserToken(token)
	assert.NoError(t, err)
	assert.Equal(t, expectedUserId, claims.UserID)
	assert.Equal(t, expectedRole, claims.Role)

}
