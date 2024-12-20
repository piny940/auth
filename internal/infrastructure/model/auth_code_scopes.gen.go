// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameAuthCodeScope = "auth_code_scopes"

// AuthCodeScope mapped from table <auth_code_scopes>
type AuthCodeScope struct {
	ScopeID    int32     `gorm:"column:scope_id;primaryKey" json:"scope_id"`
	AuthCodeID int64     `gorm:"column:auth_code_id;primaryKey" json:"auth_code_id"`
	CreatedAt  time.Time `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at;not null" json:"updated_at"`
}

// TableName AuthCodeScope's table name
func (*AuthCodeScope) TableName() string {
	return TableNameAuthCodeScope
}
