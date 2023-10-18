package main

import (
	"github.com/0xhoang/go-kit/config"
	"github.com/0xhoang/go-kit/internal/dao/daomock"
	"github.com/0xhoang/go-kit/internal/models"
	"github.com/0xhoang/go-kit/internal/must"
	"github.com/0xhoang/go-kit/internal/services"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"log"
	"testing"
)

func TestMainTestSuite(t *testing.T) {
	suite.Run(t, new(MainTestSuite))
}

type MainTestSuite struct {
	suite.Suite
	userSvc     *services.User
	mockUserDao *daomock.UserDaoMock
}

func (suite *MainTestSuite) SetupTest() {
	cfg := config.ReadConfigAndArg()
	logger, _, err := must.NewLogger(cfg.SentryDSN)
	if err != nil {
		log.Fatalf("logger: %v", err)
	}

	suite.mockUserDao = new(daomock.UserDaoMock)
	suite.userSvc = services.NewUser(logger, cfg, suite.mockUserDao)

}

func (suite *MainTestSuite) TestFindByID() {
	model := &models.User{
		FirstName: "test",
	}

	//happy case
	suite.mockUserDao.On("FindByID", uint(10)).Return(model, nil)
	user, _ := suite.userSvc.FindByID(uint(10))

	suite.Equal(user.ID, model.ID)
	suite.Equal(user.FirstName, model.FirstName)
}

func (suite *MainTestSuite) TestFindByIDError() {
	//error case
	err := errors.New("not found")
	suite.mockUserDao.On("FindByID", uint(10)).Return(&models.User{}, err)
	user1, err1 := suite.userSvc.FindByID(uint(10))

	suite.Equal(user1, (*models.User)(nil))
	suite.Equal(err, err1)
}

func (suite *MainTestSuite) TestAuthenticateByEmailPassword() {
	//email invalid
	_, err1 := suite.userSvc.AuthenticateByEmailPassword("test", "test")
	suite.Equal(err1, must.ErrInvalidEmail)

	//happy case
	email := "test@example.com"
	hashed, _ := bcrypt.GenerateFromPassword([]byte("test"), bcrypt.DefaultCost)
	user := &models.User{
		Email:    email,
		Password: string(hashed),
	}
	suite.mockUserDao.On("FindByEmail", email).Return(&models.User{
		Password: string(hashed),
		UserName: "username",
	}, nil)

	user1, err1 := suite.userSvc.AuthenticateByEmailPassword(email, "test")
	suite.EqualValues(err1, nil)
	suite.Equal(user1.Email, user.UserName)
}
