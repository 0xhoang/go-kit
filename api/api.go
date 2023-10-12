package api

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"gitlab.com/idolauncher/go-template-kit/services"
	log "go.uber.org/zap"
)

type Server struct {
	g            *gin.Engine
	authMw       *jwt.GinJWTMiddleware
	logger       *log.Logger
	helloService *services.HelloService
	userSvc      *services.User
}

func NewServer(g *gin.Engine, authMw *jwt.GinJWTMiddleware, logger *log.Logger, helloService *services.HelloService) *Server {
	return &Server{g: g, authMw: authMw, logger: logger, helloService: helloService}
}
