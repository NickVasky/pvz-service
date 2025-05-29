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

var secret []byte = []byte("the-MOST-secret-KEY-in-the-WORLD!")

type Role string

const (
	moderatorRole Role = "moderator"
	clientRole    Role = "client"
)

type endpointAccessOpts struct {
	isProtected bool
	roles       []Role
}

var endpointAccess = map[string]endpointAccessOpts{
	"DummyLoginPost": {
		isProtected: false,
		roles:       []Role{},
	},
	"RegisterPost": {
		isProtected: false,
		roles:       []Role{},
	},
	"LoginPost": {
		isProtected: false,
		roles:       []Role{},
	},
	"PvzGet": {
		isProtected: true,
		roles:       []Role{clientRole, moderatorRole},
	},
	"PvzPost": {
		isProtected: true,
		roles:       []Role{moderatorRole},
	},
	"PvzPvzIdCloseLastReceptionPost": {
		isProtected: true,
		roles:       []Role{clientRole},
	},
	"PvzPvzIdDeleteLastProductPost": {
		isProtected: true,
		roles:       []Role{clientRole},
	},
	"ReceptionsPost": {
		isProtected: true,
		roles:       []Role{clientRole},
	},
	"ProductsPost": {
		isProtected: true,
		roles:       []Role{clientRole},
	},
}

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

func AuthMiddleware(next http.Handler, allowedRoles []Role) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
			if claims.Role == string(allowedRoles[i]) {
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
