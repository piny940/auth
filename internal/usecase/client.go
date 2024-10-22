package usecase

import "auth/internal/domain/oauth"

type ClientUsecase struct {
	ClientRepo oauth.IClientRepo
}

type IClientUsecase interface {
	Find(id oauth.ClientID) (*oauth.Client, error)
}

func NewClientUsecase(clientRepo oauth.IClientRepo) IClientUsecase {
	return &ClientUsecase{
		ClientRepo: clientRepo,
	}
}

var _ IClientUsecase = &ClientUsecase{}

func (c *ClientUsecase) Find(id oauth.ClientID) (*oauth.Client, error) {
	return c.ClientRepo.FindByID(id)
}
