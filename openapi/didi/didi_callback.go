package didi

import (
	"fmt"
)

/**
滴滴状态回调通知
1、通知
2、获取订单详情
3、通知订单系统状态变更
*/
const (
	/*
	1-订单中间状态流转
	2-订单终态通知
	3-支付确认通知
	4-订单退款通知
	5-订单改价通知
	6-客服关单通知
	*/
	NT_PROCESSING   = 1
	NT_ENDSTATUS    = 2
	NT_PAY_APPLY    = 3
	NT_REFUND       = 4
	NT_CHANGE_PRICE = 5
	NT_CLOSE_BILL   = 6
)

type NotifyProcess func(request *DidiCallBackRequest) bool

type DidiCallBackRequest struct {
	Client    string `json:"client_id"`   //申请应用时分配的AppKey(同授权认证)
	Order     string `json:"order_id"`    //订单id
	Type      int    `json:"notify_type"` //通知类型
	Desc      string `json:"notify_desc"` //通知说明
	Timestamp string //当前时间戳
	Sign      string //签名

}

func (this *DidiCallBackRequest) toMap() map[string]string {
	result := make(map[string]string)
	result["client_id"] = this.Client
	result["order_id"] = this.Order
	result["notify_type"] = fmt.Sprintf("%d", this.Type)
	result["notify_desc"] = this.Desc
	result["timestamp"] = this.Timestamp
	return result

}

func (this *DiDiConfig) CallBack(p *DidiCallBackRequest) *DidiBaseResponse {
	if p != nil {
		targetSign := this.Sign(p.toMap())
		if targetSign == p.Sign {
			if this.Notify(p) {
				return &DidiBaseResponse{0, "success!"}
			}
		}
	}

	return &DidiBaseResponse{1, "process error!"}
}
