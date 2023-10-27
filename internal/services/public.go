package services

import (
	"context"
	"fmt"
	"github.com/0xhoang/go-kit/config"
	"github.com/0xhoang/go-kit/gen"
	"github.com/0xhoang/go-kit/internal/dao"
	"github.com/0xhoang/go-kit/internal/models"
	"github.com/0xhoang/go-kit/internal/must"
	"github.com/0xhoang/go-kit/internal/serializers"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"regexp"
	"time"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type goKitPublicServiceInterface interface {
	gen.GoKitPublicServer
}

var _ goKitPublicServiceInterface = (*GokitPublicService)(nil)

type GokitPublicService struct {
	logger  *zap.Logger
	cfg     *config.Config
	userDao dao.UserDaoInterface
}

func NewGokitPublicService(logger *zap.Logger, cfg *config.Config, userDao dao.UserDaoInterface) *GokitPublicService {
	return &GokitPublicService{logger: logger, cfg: cfg, userDao: userDao}
}

func (e *GokitPublicService) RegisterGrpcServer(s *grpc.Server) {
	gen.RegisterGoKitPublicServer(s, e)
}

func (e *GokitPublicService) RegisterHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	err := gen.RegisterGoKitPublicHandler(ctx, mux, conn)
	if err != nil {
		return err
	}

	return nil
}

func (g *GokitPublicService) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	return ctx, nil
}

func (e *GokitPublicService) SayHello(ctx context.Context, in *gen.HelloRequest) (*gen.HelloReply, error) {
	return &gen.HelloReply{Message: in.Name + " world"}, nil
}

func (u *GokitPublicService) isEmailValid(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	return emailRegex.MatchString(e)
}

func (u *GokitPublicService) authenticatorByEmailPassword(email, password string) (*models.User, error) {
	if !u.isEmailValid(email) {
		return nil, must.ErrInvalidEmail
	}

	user, _ := u.userDao.FindByEmail(email)
	if user != nil {
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			return nil, must.ErrInvalidPassword
		}

		return user, nil
	}

	return nil, must.ErrEmailNotExists
}

func (e *GokitPublicService) Auth(ctx context.Context, in *gen.LoginRequest) (*gen.LoginResponse, error) {
	user, err := e.authenticatorByEmailPassword(in.Email, in.Password)
	if err != nil {
		return nil, err
	}

	data := &serializers.UserInfo{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	expire := time.Now().Add(60 * time.Minute)
	token, err := must.CreateNewWithClaims(data, e.cfg.AuthenticationSecretKey, expire)

	return &gen.LoginResponse{
		Token:   token,
		Expired: fmt.Sprintf("%d", expire.Unix()),
	}, nil
}
