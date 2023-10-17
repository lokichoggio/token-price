package errorx

const (
	HttpCustomizeCodeOK = 0
)

const (
	UnknownError = iota + 10000
	ServerError
	ParamError
)

var CodeMsg = map[int]string{
	UnknownError:        "未知异常",
	ServerError:         "服务内部错误",
	ParamError:          "请求参数错误,请参考接口文档确认.",
	HttpCustomizeCodeOK: "请求成功",
}

func getMsgByCode(code int) string {
	msg, ok := CodeMsg[code]
	if !ok {
		return ""
	}
	return msg
}
