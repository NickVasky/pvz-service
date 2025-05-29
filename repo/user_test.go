package repo

import (
	cfg "AvitoTechPVZ/config"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserRepoAdd(t *testing.T) {
	cfg, err := cfg.GetConfig()
	assert.Error(t, err)
	dbconn, err := NewDbConn(cfg.Db)
	assert.Error(t, err)
	ur := UserRepo{DB: dbconn}
	roleId, _ := uuid.Parse("88e16f8e-6d96-4971-8d9a-08e40f638a19")
	newUser := User{
		Id:       uuid.New(),
		RoleId:   roleId,
		Email:    "nickvasky@gmail.com",
		Password: "testPass",
	}
	err = ur.Add(newUser)
	if err != nil {
		t.Errorf("Error is: %v", err)
	}
}

func TestUserRepoGetById(t *testing.T) {
	cfg, err := cfg.GetConfig()
	assert.Error(t, err)
	dbconn, err := NewDbConn(cfg.Db)
	assert.Error(t, err)
	ur := UserRepo{DB: dbconn}
	userId, _ := uuid.Parse("ced4ac20-7dfa-42fe-8e9a-4914e6492e8b")
	wantEmail := "be11a@yandex.ru"
	user, err := ur.GetById(userId)

	if err != nil {
		t.Errorf("Err: %v", err)
	} else if user.Email != wantEmail {
		t.Errorf("Expected user: %v, returned user: %v", wantEmail, user.Email)
	}
}
