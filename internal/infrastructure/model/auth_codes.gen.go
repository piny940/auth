// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameAuthCode = "auth_codes"

// AuthCode mapped from table <auth_codes>
type AuthCode struct {
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Value       string    `gorm:"column:value;not null" json:"value"`
	ClientID    string    `gorm:"column:client_id;not null" json:"client_id"`
	UserID      int64     `gorm:"column:user_id;not null" json:"user_id"`
	RedirectURI string    `gorm:"column:redirect_uri;not null" json:"redirect_uri"`
	Used        bool      `gorm:"column:used;not null" json:"used"`
	ExpiresAt   time.Time `gorm:"column:expires_at;not null" json:"expires_at"`
	AuthTime    time.Time `gorm:"column:auth_time;not null" json:"auth_time"`
	CreatedAt   time.Time `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;not null" json:"updated_at"`
}

// TableName AuthCode's table name
func (*AuthCode) TableName() string {
	return TableNameAuthCode
}
