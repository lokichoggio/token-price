package controller

import (
	"token-price/internal/common/errorx"
	"token-price/internal/service"

	"github.com/gin-gonic/gin"
)

func GetTokenUsdPrice(c *gin.Context) (interface{}, error) {
	token := c.Query("token")
	if token == "" {
		return nil, errorx.NewError(errorx.ParamError, "invalid token")
	}

	return service.GetTokenUsdPrice(token)
}
