package api

import (
	"github.com/0xhoang/go-kit/serializers"
	service "github.com/0xhoang/go-kit/services"
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

func (s *Server) WithAuthMw(authMw *jwt.GinJWTMiddleware) {
	s.authMw = authMw
}

func (s *Server) AuthMiddleware(key string) *jwt.GinJWTMiddleware {
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
			user, err := s.AuthenticateByEmailPassword(c)

			switch cErr := errors.Cause(err); cErr {
			case service.ErrEmailNotExists, service.ErrInactiveAccount, service.ErrInvalidPassword, service.ErrEmailIsNotVerified:
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

	return mw
}
