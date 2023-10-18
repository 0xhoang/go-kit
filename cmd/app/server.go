package main

import (
	"github.com/0xhoang/go-kit/cmd/task"
	"github.com/0xhoang/go-kit/internal/must"
	"github.com/0xhoang/go-kit/internal/services"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type Server struct {
	g      *gin.Engine
	authMw *jwt.GinJWTMiddleware
	logger *zap.Logger

	helloService *services.HelloService
	userSvc      *services.User
	eventService *task.EventService
}

func NewServer(g *gin.Engine, authMw *jwt.GinJWTMiddleware, logger *zap.Logger, helloService *services.HelloService, userSvc *services.User, eventService *task.EventService) *Server {
	return &Server{g: g, authMw: authMw, logger: logger, helloService: helloService, userSvc: userSvc, eventService: eventService}
}

func (s *Server) ErrorException(c *gin.Context, err error) {
	if err != nil {
		switch err.(type) {
		case *must.Error:
			c.JSON(http.StatusBadRequest, err.(*must.Error))
			return
		default:
			s.logger.Error("s.ErrorResponseHandle", zap.Error(err))
			c.JSON(http.StatusInternalServerError, must.ErrInternalServerError)
			return
		}
	}
}

func (s *Server) pagingFromContext(c *gin.Context) (int, int) {
	var (
		pageS  = c.DefaultQuery("page", "1")
		limitS = c.DefaultQuery("limit", "10")
		page   int
		limit  int
		err    error
	)

	page, err = strconv.Atoi(pageS)
	if err != nil {
		page = 1
	}

	limit, err = strconv.Atoi(limitS)
	if err != nil {
		limit = 10
	}

	return page, limit
}
