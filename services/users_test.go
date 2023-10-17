package services

import (
	"github.com/0xhoang/go-kit/config"
	"github.com/0xhoang/go-kit/dao/users"
	"github.com/0xhoang/go-kit/models"
	"github.com/getsentry/sentry-go"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"log"
	"testing"
	"time"
)

func TestUsersTestSuite(t *testing.T) {
	suite.Run(t, new(UsersTestSuite))
}

type UsersTestSuite struct {
	suite.Suite
	userSvc     *User
	userDaoMock *users.UserDaoMock
}

func (suite *UsersTestSuite) SetupTest() {
	mockUserDao := new(users.UserDaoMock)

	cfg := config.ReadConfigAndArg()

	sentryClient, err := sentry.NewClient(sentry.ClientOptions{
		Dsn: cfg.SentryDSN,
	})

	if err != nil {
		log.Fatal("failed to init sentry", zap.Error(err))
	}

	// Flush buffered events before the program terminates.
	// Set the timeout to the maximum duration the program can afford to wait.
	defer sentryClient.Flush(2 * time.Second)

	zapLog, _ := zap.NewProduction()
	defer zapLog.Sync()

	logger := NewLogger(zapLog, sentryClient)

	suite.userDaoMock = mockUserDao
	suite.userSvc = NewUser(logger, cfg, mockUserDao)

}

func (suite *UsersTestSuite) TestFindByID() {
	model := &models.User{
		FirstName: "test",
	}

	//happy case
	suite.userDaoMock.On("FindByID", uint(10)).Return(model, nil)
	user, _ := suite.userSvc.FindByID(uint(10))

	suite.Equal(user.ID, model.ID)
	suite.Equal(user.FirstName, model.FirstName)
}

func (suite *UsersTestSuite) TestFindByIDError() {
	//error case
	err := errors.New("not found")
	suite.userDaoMock.On("FindByID", uint(10)).Return(&models.User{}, err)
	user1, err1 := suite.userSvc.FindByID(uint(10))

	suite.Equal(user1, (*models.User)(nil))
	suite.Equal(err, err1)
}

func (suite *UsersTestSuite) TestAuthenticateByEmailPassword() {
	//email invalid
	_, err1 := suite.userSvc.AuthenticateByEmailPassword("test", "test")
	suite.Equal(err1, ErrInvalidEmail)

	//happy case
	email := "test@example.com"
	hashed, _ := bcrypt.GenerateFromPassword([]byte("test"), bcrypt.DefaultCost)
	user := &models.User{
		Email:    email,
		Password: string(hashed),
	}
	suite.userDaoMock.On("FindByEmail", email).Return(&models.User{
		Password: string(hashed),
		UserName: "username",
	}, nil)

	user1, err1 := suite.userSvc.AuthenticateByEmailPassword(email, "test")
	suite.EqualValues(err1, nil)
	suite.Equal(user1.Email, user.UserName)
}
