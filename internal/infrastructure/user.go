package infrastructure

import (
	"auth/internal/domain"
	"auth/internal/infrastructure/model"
	"auth/internal/infrastructure/query"
	"errors"

	"gorm.io/gorm"
)

type userRepo struct {
	db    *DB
	query *query.Query
}

var _ domain.IUserRepo = &userRepo{}

func NewUserRepo(db *DB) domain.IUserRepo {
	query := query.Use(db.Client)
	return &userRepo{
		db:    db,
		query: query,
	}
}

func (u *userRepo) FindByID(id int64) (*domain.User, error) {
	user, err := u.query.User.Where(u.query.User.ID.Eq(id)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrRecordNotFound{}
	}
	if err != nil {
		return nil, err
	}
	return toDomainUser(user), nil
}

func (u *userRepo) FindByName(name string) (*domain.User, error) {
	user, err := u.query.User.Where(u.query.User.Name.Eq(name)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrRecordNotFound{}
	}
	if err != nil {
		return nil, err
	}
	return toDomainUser(user), nil
}

func (u *userRepo) Create(name, encryptedPassword string) error {
	return u.query.User.Create(&model.User{
		Name:              name,
		EncryptedPassword: encryptedPassword,
	})
}

func toDomainUser(user *model.User) *domain.User {
	return &domain.User{
		ID:                domain.UserID(user.ID),
		Name:              user.Name,
		EncryptedPassword: user.EncryptedPassword,
		CreatedAt:         user.CreatedAt,
		UpdatedAt:         user.UpdatedAt,
	}
}
