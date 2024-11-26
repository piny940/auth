package gateway

import (
	"auth/internal/domain"
	"auth/internal/infrastructure"
	"auth/internal/infrastructure/query"
	"errors"
	"testing"
)

func TestUserCreate(t *testing.T) {
	setup(t)

	db := infrastructure.GetDB()
	query := query.Use(db.Client)
	userRepo := NewUserRepo(db)
	err := userRepo.Create("test@example.com", "test", "test")
	if err != nil {
		t.Fatal(err)
	}
	user, err := query.User.Where(query.User.Name.Eq("test")).First()
	if err != nil {
		t.Fatal(err)
	}
	if user.EncryptedPassword != "test" {
		t.Errorf("unexpected user.EncryptedPassword: %s", user.EncryptedPassword)
	}
	if user.Email != "test@example.com" {
		t.Errorf("unexpected user.Email: %s", user.Email)
	}
}

func TestUserFindById(t *testing.T) {
	setup(t)

	const email = "test@example.com"
	const name = "test"
	const password = "test"
	db := infrastructure.GetDB()
	query := query.Use(db.Client)
	userRepo := NewUserRepo(db)
	err := userRepo.Create(email, name, password)
	if err != nil {
		t.Fatal(err)
	}
	created, err := query.User.Where(query.User.Name.Eq(name)).First()
	if err != nil {
		t.Fatal(err)
	}
	user, err := userRepo.FindByID(domain.UserID(created.ID))
	if err != nil {
		t.Fatal(err)
	}
	if user.Name != name {
		t.Errorf("unexpected user.Name: %s", user.Name)
	}
	if user.EncryptedPassword != password {
		t.Errorf("unexpected user.EncryptedPassword: %s", user.EncryptedPassword)
	}
}

func TestUserFindByIdNotFound(t *testing.T) {
	setup(t)

	db := infrastructure.GetDB()
	userRepo := NewUserRepo(db)
	_, err := userRepo.FindByID(0)
	if !errors.Is(err, domain.ErrRecordNotFound) {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestUserFindByName(t *testing.T) {
	setup(t)

	const email = "test@example.com"
	const name = "test"
	const password = "test"

	db := infrastructure.GetDB()
	userRepo := NewUserRepo(db)
	err := userRepo.Create(email, name, password)
	if err != nil {
		t.Fatal(err)
	}
	user, err := userRepo.FindByName(name)
	if err != nil {
		t.Fatal(err)
	}
	if user.Name != name {
		t.Errorf("unexpected user.Name: %s", user.Name)
	}
	if user.EncryptedPassword != password {
		t.Errorf("unexpected user.EncryptedPassword: %s", user.EncryptedPassword)
	}
}

func TestUserFindByNameNotFound(t *testing.T) {
	setup(t)

	db := infrastructure.GetDB()
	userRepo := NewUserRepo(db)
	_, err := userRepo.FindByName("test")
	if !errors.Is(err, domain.ErrRecordNotFound) {
		t.Errorf("unexpected error: %v", err)
	}
}
