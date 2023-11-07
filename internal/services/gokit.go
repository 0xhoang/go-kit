package services

import (
	"context"
	"github.com/0xhoang/go-kit/common"
	"github.com/0xhoang/go-kit/config"
	"github.com/0xhoang/go-kit/gen"
	"github.com/0xhoang/go-kit/internal/dao"
	"github.com/0xhoang/go-kit/internal/models"
	"github.com/0xhoang/go-kit/internal/must"
	"github.com/0xhoang/go-kit/internal/serializers"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type goKitServiceInterface interface {
	gen.GoKitServer
}

var _ goKitServiceInterface = (*GokitService)(nil)

type GokitService struct {
	logger  *zap.Logger
	cfg     *config.Config
	db      *gorm.DB
	userDao dao.UserDaoInterface
}

func NewGokitService(logger *zap.Logger, cfg *config.Config, db *gorm.DB, userDao dao.UserDaoInterface) *GokitService {
	return &GokitService{logger: logger, cfg: cfg, db: db, userDao: userDao}
}

func (e *GokitService) RegisterGrpcServer(s *grpc.Server) {
	gen.RegisterGoKitServer(s, e)
}

func (e *GokitService) RegisterHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	err := gen.RegisterGoKitHandler(ctx, mux, conn)
	if err != nil {
		return err
	}

	return nil
}

func (s *GokitService) userFromContext(c context.Context) (*models.User, error) {
	authen := c.Value(common.CustomerKey)
	if authen == nil {
		s.logger.Error("authen empty")
		return nil, must.ErrInvalidCredentials
	}

	userIDVal := authen.(*serializers.UserInfo)

	user, err := s.userDao.FindByID(userIDVal.ID)
	if err != nil || user == nil {
		s.logger.Error("FindByID", zap.Error(err))
		return nil, must.ErrInvalidCredentials
	}

	return user, nil
}
