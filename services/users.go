package services

import (
	"github.com/0xhoang/go-kit/config"
	"github.com/0xhoang/go-kit/dao/users"
	"github.com/0xhoang/go-kit/models"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"regexp"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type User struct {
	logger *zap.Logger
	conf   *config.Config
	r      users.UserDaoInterface
}

func NewUser(logger *zap.Logger, conf *config.Config, r users.UserDaoInterface) *User {
	return &User{logger: logger, conf: conf, r: r}
}

func (u *User) FindByID(id uint) (*models.User, error) {
	user, err := u.r.FindByID(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *User) isEmailValid(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	return emailRegex.MatchString(e)
}

func (u *User) AuthenticateByEmailPassword(email, password string) (*models.User, error) {
	if !u.isEmailValid(email) {
		return nil, ErrInvalidEmail
	}

	user, _ := u.r.FindByEmail(email)

	if user != nil {
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			return nil, ErrInvalidPassword
		}

		return user, nil
	}

	return nil, ErrEmailNotExists
}
