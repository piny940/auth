package infrastructure

import (
	"auth/internal/domain"
	"auth/internal/infrastructure/query"
)

type userRepo struct {
	db    *DB
	query *query.Query
}

// FindByID implements domain.IUserRepo.
func (u *userRepo) FindByID(id int64) (*domain.User, error) {
	panic("unimplemented")
}

// FindByName implements domain.IUserRepo.
func (u *userRepo) FindByName(name string) (*domain.User, error) {
	panic("unimplemented")
}

var _ domain.IUserRepo = &userRepo{}

func NewUserRepo(db *DB) domain.IUserRepo {
	query := query.Use(db.Client)
	return &userRepo{
		db:    db,
		query: query,
	}
}
