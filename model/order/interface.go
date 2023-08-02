package order

type Order struct {
	Id     int32 `json:"id"`
	UserId int32 `json:"user_id"`
	ItemId int32 `json:"item_id"`
	Status int32 `json:"status"` // 1：已创建，等待支付 2：已超时，订单失效 3：手动取消 4：已支付
}
