package api

import (
	"github.com/0xhoang/go-kit/services"
	"github.com/0xhoang/go-kit/task"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	log "go.uber.org/zap"
)

type Server struct {
	g            *gin.Engine
	authMw       *jwt.GinJWTMiddleware
	logger       *log.Logger
	helloService *services.HelloService
	userSvc      *services.User
	eventService *task.EventService
}

func NewServer(g *gin.Engine, authMw *jwt.GinJWTMiddleware, logger *log.Logger, helloService *services.HelloService, userSvc *services.User, eventService *task.EventService) *Server {
	return &Server{g: g, authMw: authMw, logger: logger, helloService: helloService, userSvc: userSvc, eventService: eventService}
}
