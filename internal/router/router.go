package router

import (
	"net/http"

	"token-price/internal/config"
	"token-price/internal/controller"
	"token-price/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Router(c *config.Config) *gin.Engine {
	gin.SetMode(c.GinMode)
	e := gin.New()

	// middlewares
	e.Use(middleware.Recovery())
	e.Use(middleware.PromMiddleware(e, nil))
	middleware.Register(e)

	// ping
	e.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	v1 := e.Group("/api/v1")
	v1.GET("/get_token_usd_price", middleware.WithJSONResp(controller.GetTokenUsdPrice))

	return e
}
