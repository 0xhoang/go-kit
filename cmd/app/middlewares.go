package main

import (
	"github.com/0xhoang/go-kit/internal/models"
	"github.com/0xhoang/go-kit/internal/must"
	"github.com/0xhoang/go-kit/internal/serializers"
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

const (
	userIDKey     = "id"
	userEmailKey  = "email"
	userFirstName = "name"
)

func (s *Server) AuthMiddleware(key string) {
	mw, _ := jwt.New(&jwt.GinJWTMiddleware{
		Key:         []byte(key),
		Timeout:     time.Hour * 24 * 7, //7 days
		MaxRefresh:  time.Hour * 24 * 7, //7 days
		IdentityKey: userIDKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*serializers.UserByEmailResp); ok {
				return jwt.MapClaims{
					userIDKey:     v.ID,
					userEmailKey:  v.Email,
					userFirstName: v.FirstName,
				}
			}

			return jwt.MapClaims{}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			user, err := s.authenticateByEmailPassword(c)

			switch cErr := errors.Cause(err); cErr {
			case must.ErrEmailNotExists,
				must.ErrInactiveAccount,
				must.ErrInvalidPassword,
				must.ErrEmailIsNotVerified:
				return nil, cErr
			case nil:
				return user, nil
			default:
				return nil, err
			}
		},
		HTTPStatusMessageFunc: func(err error, c *gin.Context) string {
			c.Set("authorize_error", err)
			return err.Error()
		},
		Unauthorized: func(c *gin.Context, _ int, _ string) {
			err, _ := c.Get("authorize_error")
			c.JSON(http.StatusUnauthorized, serializers.Resp{
				Error: err,
			})
		},
		LoginResponse: func(c *gin.Context, code int, token string, expire time.Time) {
			c.JSON(http.StatusOK, serializers.Resp{
				Result: serializers.UserLoginResp{
					Token:   token,
					Expired: expire.Format(time.RFC3339),
				},
				Error: nil,
			})
		},
		RefreshResponse: func(c *gin.Context, code int, token string, expire time.Time) {
			c.JSON(http.StatusOK, serializers.Resp{
				Result: serializers.UserLoginResp{
					Token:   token,
					Expired: expire.Format(time.RFC3339),
				},
				Error: nil,
			})
		},
	})

	s.authMw = mw
}

func (s *Server) authenticateByEmailPassword(c *gin.Context) (*serializers.UserByEmailResp, error) {
	var req serializers.AuthByEmailReq
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, must.ErrInvalidArgument
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
