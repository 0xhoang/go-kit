package dao

import "github.com/0xhoang/go-kit/internal/models"

type UserDaoInterface interface {
	FindByID(id uint) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
}
