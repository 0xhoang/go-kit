package api

import (
	"github.com/0xhoang/go-kit/serializers"
	"github.com/0xhoang/go-kit/services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
		s.logger.Error("s.apiService.HelloWorld", zap.Error(err))
		c.JSON(http.StatusInternalServerError, serializers.Resp{Error: services.ErrInternalServerError})
		return
	}

	c.JSON(http.StatusOK, serializers.Resp{Result: resp})
}
