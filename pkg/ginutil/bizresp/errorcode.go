package bizresp

/*
错误码统一为5位整数，前三位按照http状态码进行分组，后两位可根据业务细分从00到99。
http状态码参考： https://www.iana.org/assignments/http-status-codes/http-status-codes.xhtml
*/

var ( // 400xx BadRequest
	InvalidParam        = errcode{40000, "参数错误"}  // 参数基本校验失败
	InvalidAssociatedID = errcode{40001, "无效的参数"} // 提交数据关联的ID无效
)

var ( // 401xx Unauthorized
	Unauthorized    = errcode{40100, "登录校验失效"}
	CaptchaExpired  = errcode{40101, "验证码过期"}
	CaptchaWrong    = errcode{40102, "验证码错误"}
	UserOrPwdWrong  = errcode{40103, "用户名或密码错误"}
	PasswordExpired = errcode{40104, "密码已过期"}
	AccountDisabled = errcode{40105, "账号不可用"}
)

var ( // 403xx Forbidden
	Forbidden      = errcode{40300, "禁止访问"}
	PermissionDeny = errcode{40301, "未分配该权限"}
	OperationDeny  = errcode{40302, "禁止操作该数据"}
)

var ( // 404xx NotFound
	RecordNotFound = errcode{40400, "记录未找到"}
)

var ( // 409xx Conflict
	Conflict = errcode{40900, "提交冲突"} // 多人操作同一资源时，后提交者版本校验不通过
)

var ( // 413 RequestEntityTooLarge
	EntityTooLarge = errcode{41300, "提交数据过大"}
)

var ( // 415 UnsupportedMediaType
	UnsupportedMediaType = errcode{41500, "不支持的媒体类型"}
)

var ( // 422xx UnprocessableEntity
	UniqueKeyExist     = errcode{42201, "唯一标识已存在"} // 唯一标识已存在
	NeedModified       = errcode{42202, "数据未修改"}   // 数据应当修改而未修改
	CodeExpiredOrWrong = errcode{42203, "临时凭证过期或已失效"}
)

var ( // 423xx Locked
	Locked = errcode{42300, "资源被锁定"} // 数据处于正在处理的中间状态
)

var ( // 429xx TooManyRequests
	TooManyRequests = errcode{42900, "请求过于频繁"} // 短时间内重复发起请求
)

var ( // 500xx InternalServerError
	InternalServerError = errcode{50000, "系统繁忙"} // 服务器内部错误
	ServerCommonError   = errcode{50001, "系统繁忙"} // 服务端通用错误
	ServerRedisError    = errcode{50002, "系统繁忙"}
	ServerNsqError      = errcode{50003, "系统繁忙"}
)

var ( // 502xx BadGateway
	ResponseWrong = errcode{50200, "网关响应错误"} // 外部接口响应错误
)

var ( // 503xx ServiceUnavailable
	ServiceUnavailable = errcode{50300, "服务暂不可用"} // 停服升级维护中
)

var ( // 504xx GatewayTimeout
	ResponseTimeout = errcode{50400, "网关响应超时"} // 外部接口请求超时
)
