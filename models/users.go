package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName  string
	MiddleName string
	LastName   string
	FullName   string
	UserName   string
	Email      string `gorm:"uniqueIndex:idx_email;size:50"`
	Password   string
	Bio        string

	IsActive        bool `gorm:"index:idx_is_active"`
	IsVerifiedEmail bool
}
