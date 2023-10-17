package middleware

import (
	"net/http"

	"token-price/internal/common/errorx"

	"github.com/gin-gonic/gin"
)

func WithJSONResp(f func(c *gin.Context) (interface{}, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := f(c)
		if data == nil {
			data = make(map[string]interface{})
		}

		switch e := err.(type) {
		case nil:
			c.JSON(http.StatusOK, Response{
				Code: errorx.HttpCustomizeCodeOK,
				Msg:  "success",
				Data: data,
			})
		case *errorx.BaseError:
			c.JSON(http.StatusOK, Response{
				Code: e.Code(),
				Msg:  e.Error(),
				Data: data,
			})

		default:
			switch err {
			//case sql.ErrNoRows:
			default:
				serverError := errorx.NewError(errorx.ServerError, err.Error())
				c.JSON(http.StatusOK, Response{
					Code: serverError.Code(),
					Msg:  serverError.Error(),
					Data: data,
				})
			}
		}
	}
}
