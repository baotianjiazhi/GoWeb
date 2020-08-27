package serializer

type Rescode int64

const (
	CodeSuccess  Rescode = 1000 + iota
	CodeInvalidParam
	CodeUserExist
	CodeUserNotExist
	CodeInvalidPassword
	CodeCreateUserFault
	CodeServerBusy
	CodeNeedLogin
	CodeInvalidToken
)



var CodeMsgMap = map[Rescode]string{
	CodeSuccess: "success",
	CodeInvalidParam: "请求参数错误",
	CodeUserExist: "用户名存在",
	CodeUserNotExist: "用户名不存在",
	CodeInvalidPassword: "用户名或密码错误",
	CodeCreateUserFault: "创建用户失败",
	CodeServerBusy: "服务器繁忙",
	CodeNeedLogin: "请先登陆",
	CodeInvalidToken: "无效的Token",
}

func (code Rescode)getMsg() (string) {
	msg, ok := CodeMsgMap[code]
	if !ok {
		msg = CodeMsgMap[CodeServerBusy]
	}

	return msg
}