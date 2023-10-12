package services

import (
	"github.com/pkg/errors"
	"gitlab.com/idolauncher/go-template-kit/config"
	"gitlab.com/idolauncher/go-template-kit/dao"
	"gitlab.com/idolauncher/go-template-kit/models"
	"go.uber.org/zap"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type User struct {
	conf   *config.Config
	logger *zap.Logger
	r      *dao.User
}

func NewUser(conf *config.Config, logger *zap.Logger, r *dao.User) *User {
	return &User{conf: conf, logger: logger, r: r}
}

func (u *User) FindByID(id uint) (*models.User, error) {
	user, err := u.r.FindByID(id)
	if err != nil {
		return nil, errors.Wrap(err, "u.portalDao.FindByID")
	}
	return user, nil
}

func (u *User) AuthenticateByEmailPassword(email, password string) (*models.User, error) {
	return nil, nil
}
