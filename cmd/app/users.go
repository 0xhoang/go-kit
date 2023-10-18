package main

import (
	"github.com/0xhoang/go-kit/internal/must"
	"github.com/0xhoang/go-kit/internal/serializers"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

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
		c.JSON(http.StatusUnauthorized, serializers.Resp{Error: must.ErrInternalServerError})
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
