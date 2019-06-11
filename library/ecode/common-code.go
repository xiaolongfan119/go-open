package ecode

var (
	OK          = Add(0, "")
	ServerErr   = Add(500, "服务器错误")
	Unknown     = Add(5000, "未知错误")
	UnknownCode = Add(5001, "未知错误码")

	TokenEmpty   = Add(2000, "token为空")
	TokenExpired = Add(2001, "token已过期")
	TokenInvalid = Add(2002, "token 无效")

	SqlInvalid = Add(3000, "sql错误")

	ParamsFormatError_1  = Add(4000, "参数格式有误(map)")
	ParamsFormatError_2  = Add(4001, "参数格式有误(validate)")
	ParamsInValid        = Add(4002, "参数无效")
	ParamsLogin          = Add(4003, "请输入账号或密码")
	AccountOrPasswordErr = Add(4004, "账号或密码错误")
	LoginNameExist       = Add(4005, "账号已存在")

	BreakerTooManyRequests = Add(6000, "breaker too many requests")
	BreakerOpenState       = Add(6001, "breaker open state")
)
