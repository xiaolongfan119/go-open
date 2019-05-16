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

	ParamsFormatError = Add(4000, "参数格式有误")
)
