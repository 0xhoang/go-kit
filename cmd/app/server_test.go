package main

import (
	"context"
	"github.com/0xhoang/go-kit/config"
	"github.com/0xhoang/go-kit/gen"
	"github.com/0xhoang/go-kit/internal/dao/daomock"
	"github.com/0xhoang/go-kit/internal/models"
	"github.com/0xhoang/go-kit/internal/must"
	"github.com/0xhoang/go-kit/internal/services"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"log"
	"testing"
)

type MainTestSuite struct {
	suite.Suite
	svc         *services.GokitService
	svcPublic   *services.GokitPublicService
	mockUserDao *daomock.UserDaoMock
}

func (suite *MainTestSuite) SetupTest() {
	cfg := config.ReadConfigAndArg()
	logger, _, err := must.NewLogger(cfg.SentryDSN, "testing")
	if err != nil {
		log.Fatalf("logger: %v", err)
	}

	suite.mockUserDao = new(daomock.UserDaoMock)
	suite.svc = services.NewGokitService(logger, cfg, nil, suite.mockUserDao)
	suite.svcPublic = services.NewGokitPublicService(logger, cfg, suite.mockUserDao)

}

func (suite *MainTestSuite) TestAuthenticateByEmailPassword() {
	//email invalid
	_, err1 := suite.svcPublic.Auth(context.Background(), &gen.LoginRequest{Email: "test", Password: "test"})
	suite.Equal(err1, must.ErrInvalidEmail)

	//happy case
	email := "test@example.com"
	hashed, _ := bcrypt.GenerateFromPassword([]byte("test"), bcrypt.DefaultCost)

	/*	user := &models.User{
		Email:    email,
		Password: string(hashed),
	}*/

	suite.mockUserDao.On("FindByEmail", email).Return(&models.User{
		Password: string(hashed),
		UserName: "username",
	}, nil)

	user1, err1 := suite.svcPublic.Auth(context.Background(), &gen.LoginRequest{Email: email, Password: "test"})

	suite.EqualValues(err1, nil)
	suite.NotEmpty(user1.Expired)
}

func TestMainTestSuite(t *testing.T) {
	suite.Run(t, new(MainTestSuite))
}
