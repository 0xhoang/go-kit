package main

import (
	"github.com/0xhoang/go-kit/internal/serializers"
	"github.com/gin-gonic/gin"
	"net/http"
)

// HelloWorld godoc
// @Description  HelloWorld
// @Tags         HelloWorld
// @Success      200  {object}  serializers.Resp
// @Failure      400  {object}  serializers.Resp
// @Router       /hello [get]
func (s *Server) HelloWorld(c *gin.Context) {
	resp, err := s.helloService.HelloWorld()

	if err != nil {
		s.ErrorException(c, err)
		return
	}

	c.JSON(http.StatusOK, serializers.Resp{Result: resp})
}
