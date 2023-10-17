package errorx

import "fmt"

type BaseError struct {
	code int
	msg  string
}

func NewError(code int, errMsg string) *BaseError {
	codeMsg := getMsgByCode(code)

	msg := codeMsg
	if errMsg != "" {
		msg = fmt.Sprintf("[%s]:%s", codeMsg, errMsg)
	}
	return &BaseError{code: code, msg: msg}
}

func (e *BaseError) Code() int {
	if e == nil {
		return HttpCustomizeCodeOK
	}
	return e.code
}

func (e *BaseError) Error() string {
	if e == nil {
		return ""
	}
	return e.msg
}
