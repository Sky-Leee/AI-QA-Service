package consts

type Code uint

const (
	CodeSuccess Code = iota + 1000
	CodeInternalErr
	CodeServerBusy
	CodeInvalidParam
	CodeNotFound
	CodeUnsupportedAuthProtocol
	CodeInvalidToken
	CodeExpiredToken

	CodeUserExist
	CodeUserNotExist
	CodeWrongPassword
	CodeNeedLogin
	CodeExpiredLogin

	CodeForbidden
	CodeTimeOut
)

var codeMsgMap = map[Code]string{
	CodeSuccess:                 "成功",
	CodeInternalErr:             "服务繁忙",
	CodeServerBusy:              "触发限流",
	CodeInvalidParam:            "无效参数",
	CodeNotFound:                "未找到",
	CodeUnsupportedAuthProtocol: "不支持的认证协议",
	CodeInvalidToken:            "无效 Token",
	CodeExpiredToken:            "过期 Token",

	CodeUserExist:     "用户已存在",
	CodeUserNotExist:  "用户不存在",
	CodeWrongPassword: "密码错误",
	CodeNeedLogin:     "需要登录",
	CodeExpiredLogin:  "登录过期",

	CodeForbidden: "禁止访问",
	CodeTimeOut:   "请求超时",
}

func (c Code) GetMsg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		return "无效错误码"
	}
	return msg
}

const (
	QUESTION_BATCH_SIZE        = 10    // 一次获取 10 道题
	QUESTION_BANK_MAX_SIZE     = 10000 // 题库最大容量
	QUESTION_LIST_NOT_FINISHED = 0     // 题单未完成
	QUESTION_LIST_FINISHED     = 1     // 题单已完成
)

const (
	COZE_GEN_KEYWORD   = "生成题单"
	COZE_JUDGE_KEYWORD = `题目标题: %v
题目内容:
%v
用户解答: %v
`
)
