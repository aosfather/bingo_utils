package fengzheng

import (
	"fmt"
	utils "github.com/aosfather/bingo_utils"
)

/*
  异步订购流量产品接口
*/

type orderResponse struct {
	Id      string `json:"orderid"`
	TradeNo string `json:"tradeno"`
	Status  string `json:"status"`

	Msg string `json:"message"`
}

type orderRequest struct {
	baseParames
	Phone     string `json:"phone"`
	ProductId string `json:"productid"`
	Time      string `json:"time"`
	TradeNo   string `json:"tradno"`
}

func (this *orderRequest) ToMap() map[string]string {
	result := make(map[string]string)
	result["phone"] = this.Phone
	result["productid"] = this.ProductId
	result["time"] = this.Time
	result["tradno"] = this.TradeNo
	return result
}

//流量充值异步
func (this *Fengzheng) Order(no, phone string, pid string) {
	this._order("order", no, phone, pid)
}

//流量充值同步实时
func (this *Fengzheng) RealTimeOrder(no, phone string, pid string) {
	this._order("timely/order", no, phone, pid)
}

//话费充值异步
func (this *Fengzheng) CallOrder(no, phone string, pid string) {
	this._order("callorder", no, phone, pid)
}

func (this *Fengzheng) _order(method string, no, phone string, pid string) *orderResponse {
	result := orderResponse{}
	request := orderRequest{}
	request.TradeNo = no
	request.Phone = phone
	request.ProductId = pid
	request.Time = utils.Now()

	this.Call(method, &request, &result)

	fmt.Println(result)
	return &result
}
