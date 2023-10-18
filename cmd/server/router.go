package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (s *Server) Routes() {
	s.g.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	s.g.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	s.g.GET("/hello", s.HelloWorld)

	auth := s.g.Group("/auth")
	auth.POST("/login", s.Login)

	me := s.g.Group("/me")
	me.Use(s.authMw.MiddlewareFunc())
	{
		me.GET("", s.Profile)
	}
}
