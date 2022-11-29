package controllers

type ResCode int64

const (
	CodeSuccess ResCode = 1000 + iota
	CodeInvalidParams
	CodeSeverBusy
	CodeNotLogin
	CodeInvalidToken
)

var codeMsgMap = map[ResCode]string{
	CodeSuccess:       "success",
	CodeSeverBusy:     "服务器繁忙",
	CodeInvalidParams: "请求参数异常",
}

func (c ResCode) String() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		msg = codeMsgMap[CodeSeverBusy]
	}
	return msg
}
