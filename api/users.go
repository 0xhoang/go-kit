package api

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gitlab.com/idolauncher/go-template-kit/models"
	"gitlab.com/idolauncher/go-template-kit/serializers"
	service "gitlab.com/idolauncher/go-template-kit/services"
)

func (s *Server) AuthenticateByEmailPassword(c *gin.Context) (*models.User, error) {
	var req serializers.AuthByEmailReq
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, service.ErrInvalidArgument
	}

	return s.userSvc.AuthenticateByEmailPassword(req.Email, req.Password)
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
