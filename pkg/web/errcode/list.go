package errcode

/*
错误码统一为5位整数，前三位按照http状态码进行分组，后两位可根据业务细分从00到99。
http状态码参考： https://www.iana.org/assignments/http-status-codes/http-status-codes.xhtml
*/

var ( // 400xx BadRequest
	InvalidParam        = ErrCode{40000, "参数错误"}  // 参数基本校验失败
	InvalidAssociatedID = ErrCode{40001, "无效的参数"} // 提交数据关联的ID无效
)

var ( // 401xx Unauthorized
	Unauthorized    = ErrCode{40100, "登录校验失效"}
	CaptchaExpired  = ErrCode{40101, "验证码过期"}
	CaptchaWrong    = ErrCode{40102, "验证码错误"}
	UserOrPwdWrong  = ErrCode{40103, "用户名或密码错误"}
	PasswordExpired = ErrCode{40104, "密码已过期"}
	AccountDisabled = ErrCode{40105, "账号不可用"}
)

var ( // 403xx Forbidden
	Forbidden      = ErrCode{40300, "禁止访问"}
	PermissionDeny = ErrCode{40301, "未分配该权限"}
	OperationDeny  = ErrCode{40302, "禁止操作该数据"}
)

var ( // 404xx NotFound
	RecordNotFound = ErrCode{40400, "记录未找到"}
)

var ( // 409xx Conflict
	Conflict = ErrCode{40900, "提交冲突"} // 多人操作同一资源时，后提交者版本校验不通过
)

var ( // 413 RequestEntityTooLarge
	EntityTooLarge = ErrCode{41300, "提交数据过大"}
)

var ( // 415 UnsupportedMediaType
	UnsupportedMediaType = ErrCode{41500, "不支持的媒体类型"}
)

var ( // 422xx UnprocessableEntity
	UniqueKeyExist     = ErrCode{42201, "唯一标识已存在"} // 唯一标识已存在
	NeedModified       = ErrCode{42202, "数据未修改"}   // 数据应当修改而未修改
	CodeExpiredOrWrong = ErrCode{42203, "临时凭证过期或已失效"}
)

var ( // 423xx Locked
	Locked = ErrCode{42300, "资源被锁定"} // 数据处于正在处理的中间状态
)

var ( // 429xx TooManyRequests
	TooManyRequests = ErrCode{42900, "请求过于频繁"} // 短时间内重复发起请求
)

var ( // 500xx InternalServerError
	InternalServerError = ErrCode{50000, "系统繁忙"} // 服务器内部错误
	ServerCommonError   = ErrCode{50001, "系统繁忙"} // 服务端通用错误
	ServerRedisError    = ErrCode{50002, "系统繁忙"}
	ServerNsqError      = ErrCode{50003, "系统繁忙"}
)

var ( // 502xx BadGateway
	ResponseWrong = ErrCode{50200, "网关响应错误"} // 外部接口响应错误
)

var ( // 503xx ServiceUnavailable
	ServiceUnavailable = ErrCode{50300, "服务暂不可用"} // 停服升级维护中
)

var ( // 504xx GatewayTimeout
	ResponseTimeout = ErrCode{50400, "网关响应超时"} // 外部接口请求超时
)
