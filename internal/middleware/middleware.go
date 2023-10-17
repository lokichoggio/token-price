package middleware

import "token-price/internal/common/errorx"

const (
	defaultErrCode = errorx.ServerError
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}
