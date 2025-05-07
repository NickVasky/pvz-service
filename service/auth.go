package service

import (
	"AvitoTechPVZ/repo"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

//const secret string = "I'm the KEY!"

var secret []byte = []byte("the-MOST-secret-KEY-in-the-WORLD!")

type UserClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func NewUserClaims(user repo.User, role repo.Role) UserClaims {
	return UserClaims{
		UserID: user.Id.String(),
		Role:   role.Name,
	}
}

func NewDummyUser(roleName string) UserClaims {
	user := repo.User{
		Id: uuid.New(),
	}
	role := repo.Role{
		Name: roleName,
	}
	return NewUserClaims(user, role)
}

func (c UserClaims) GenerateJwtToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	signedToken, err := token.SignedString(secret)
	if err != nil {
		panic(err)
	}
	fmt.Println(signedToken) //DBG
	return signedToken
}

func ParseUserToken(token string) (*UserClaims, error) {
	validMethods := []string{"HS256"}
	claims := &UserClaims{}

	parsedToken, err := jwt.ParseWithClaims(token, claims, jwtKeyFunc, jwt.WithValidMethods(validMethods))

	if err != nil {
		return nil, err
	}
	if !parsedToken.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil

}

func jwtKeyFunc(token *jwt.Token) (interface{}, error) {
	// Ensure the signing method is correct
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, errors.New("unexpected signing method")
	}
	return secret, nil
}

func AuthMiddleware(next http.HandlerFunc, allowedRoles []string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Example: Check for a valid Authorization header
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		authHeaderParts := strings.Split(authHeader, " ")
		if !(len(authHeaderParts) == 2 && authHeaderParts[0] == "Bearer") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		token := authHeaderParts[1]
		claims, err := ParseUserToken(token)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		isAllowed := false
		for i := range allowedRoles {
			if claims.Role == allowedRoles[i] {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
