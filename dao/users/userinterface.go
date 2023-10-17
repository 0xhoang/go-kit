package users

import "github.com/0xhoang/go-kit/models"

type UserDaoInterface interface {
	FindByID(id uint) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
}
