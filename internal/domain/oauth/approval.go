package oauth

import (
	"auth/internal/domain"
	"time"
)

type ApprovalID int64
type Approval struct {
	ID        ApprovalID
	ClientID  ClientID
	UserID    domain.UserID
	Scopes    []TypeScope
	CreatedAt time.Time
	UpdatedAt time.Time
}
type IApprovalRepo interface {
	Find(clientID ClientID, userID domain.UserID) (*Approval, error)
	Create(clientID ClientID, userID domain.UserID, scopes []TypeScope) error
}
