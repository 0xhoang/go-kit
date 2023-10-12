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
	Email      string
	Password   string
	Bio        string `gorm:"type:text"`

	IsActive        bool
	IsVerifiedEmail bool
}
