package api

import (
	"github.com/0xhoang/go-kit/models"
	"github.com/0xhoang/go-kit/serializers"
	service "github.com/0xhoang/go-kit/services"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
)

func (s *Server) AuthenticateByEmailPassword(c *gin.Context) (*serializers.UserByEmailResp, error) {
	var req serializers.AuthByEmailReq
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, service.ErrInvalidArgument
	}

	user, err := s.userSvc.AuthenticateByEmailPassword(req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	resp := &serializers.UserByEmailResp{
		ID:              user.ID,
		Email:           user.Email,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		IsActive:        user.IsActive,
		IsVerifiedEmail: user.IsVerifiedEmail,
	}

	return resp, nil
}

func (s *Server) userFromContext(c *gin.Context) (*models.User, error) {
	userIDVal, ok := c.Get(userIDKey)
	if !ok {
		return nil, errors.New("failed to get userIDKey from context")
	}

	userID := userIDVal.(float64)
	user, err := s.userSvc.FindByID(uint(userID))
	if err != nil {
		return nil, errors.Wrap(err, "s.userSvc.FindByID")
	}
	return user, nil
}

// Login godoc
// @Tags         Auth
// @Param        request body serializers.AuthByEmailReq  true  "payload"
// @Success      200  {object}  serializers.Resp
// @Failure      400  {object}  services.Error
// @Failure      500  {object}  services.Error
// @Router       /auth/login [post]
func (s *Server) Login(c *gin.Context) {
	s.authMw.LoginHandler(c)
}

// Profile godoc
// @Tags         Auth
// @Param		 Authorization  header  string  true  "Bearer Token"
// @Success      200  {object}  serializers.Resp
// @Failure      400  {object}  services.Error
// @Failure      500  {object}  services.Error
// @Router       /me [get]
func (s *Server) Profile(c *gin.Context) {
	user, err := s.userFromContext(c)
	if err != nil {
		s.logger.Error("s.userFromContext", zap.Error(err))
		c.JSON(http.StatusUnauthorized, serializers.Resp{Error: service.ErrInternalServerError})
		return
	}

	resp := &serializers.UserByEmailResp{
		ID:              user.ID,
		Email:           user.Email,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		IsActive:        user.IsActive,
		IsVerifiedEmail: user.IsVerifiedEmail,
	}

	c.JSON(http.StatusOK, serializers.Resp{Result: resp})
}
