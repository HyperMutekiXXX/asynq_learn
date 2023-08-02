package constant

const (
	RedisAddr    = "120.79.13.57:6379"
	OrderMission = "order:status"
)

const (
	PayWait    = 1 // 等待支付
	PayTimeOut = 2 // 支付超时
	PayCancel  = 3 // 取消支付
	PayOk      = 4 // 已支付
)
