package repo

import (
	"testing"

	"github.com/google/uuid"
)

func TestUserRepoAdd(t *testing.T) {
	ur := UserRepo{DB: OpenDbConnection()}
	roleId, _ := uuid.Parse("88e16f8e-6d96-4971-8d9a-08e40f638a19")
	newUser := User{
		Id:       uuid.New(),
		RoleId:   roleId,
		Email:    "be11a@yandex.ru",
		Password: "thebestgirl",
	}
	err := ur.Add(newUser)
	if err != nil {
		t.Errorf("Error is: %v", err)
	}
}

func TestUserRepoGetById(t *testing.T) {
	ur := UserRepo{DB: OpenDbConnection()}
	userId, _ := uuid.Parse("ced4ac20-7dfa-42fe-8e9a-4914e6492e8b")
	wantEmail := "be11a@yandex.ru"
	user, err := ur.GetById(userId)

	if err != nil {
		t.Errorf("Err: %v", err)
	} else if user.Email != wantEmail {
		t.Errorf("Expected user: %v, returned user: %v", wantEmail, user.Email)
	}
}
