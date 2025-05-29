package security

import (
	"AvitoTechPVZ/repo"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// TODO: store secret in env?
//var secret []byte = []byte("the-MOST-secret-KEY-in-the-WORLD!")

type Role string

const (
	moderatorRole Role = "moderator"
	clientRole    Role = "client"
)

type endpointAccessRule struct {
	isProtected bool
	roles       []Role
}

var endpointAccessRules = map[string]endpointAccessRule{
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

type SecurityController struct {
	secret []byte
}

func NewSecurityController(key string) *SecurityController {
	return &SecurityController{secret: []byte(key)}
}

func (s *SecurityController) getSecret() []byte {
	log.Printf("Secret requested: %v", s.secret)
	return s.secret
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

func (s *SecurityController) GenerateJwtToken(c UserClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	signedToken, err := token.SignedString(s.getSecret())
	if err != nil {
		panic(err)
	}
	fmt.Println(signedToken) //DBG
	return signedToken
}

func (s *SecurityController) ParseUserToken(token string) (*UserClaims, error) {
	validMethods := []string{"HS256"}
	claims := &UserClaims{}

	parsedToken, err := jwt.ParseWithClaims(token, claims, s.jwtKeyFunc, jwt.WithValidMethods(validMethods))

	if err != nil {
		return nil, err
	}
	if !parsedToken.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil

}

func (s *SecurityController) jwtKeyFunc(token *jwt.Token) (interface{}, error) {
	// Ensure the signing method is correct
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, errors.New("unexpected signing method")
	}
	return s.getSecret(), nil
}

func (s *SecurityController) AuthMiddleware(inner http.Handler, endpointName string) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		endpointAccess, ok := endpointAccessRules[endpointName]
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		// skip auth
		if !endpointAccess.isProtected {
			inner.ServeHTTP(w, r)
		}

		allowedRoles := endpointAccess.roles
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
		claims, err := s.ParseUserToken(token)
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

		inner.ServeHTTP(w, r)
	})
}
