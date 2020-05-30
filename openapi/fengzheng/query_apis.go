package fengzheng

import (
	"fmt"
	utils "github.com/aosfather/bingo_utils"
)

/**
  查询接口
   1、查询订单状态
   2、余额查询
*/

type queryRequest struct {
	baseParames
	TradeNo string `json:"tradeno"`
	OrderId int64  `json:"orderId"`
}

func (this *queryRequest) ToMap() map[string]string {
	result := make(map[string]string)
	result["orderId"] = fmt.Sprintf("%d", this.OrderId)
	result["tradeno"] = this.TradeNo
	return result
}

type queryResponse struct {
	baseResponse
	QueryStatus string `json:"queryStatus"` //查询状态码
	TradeNo     string `json:"tradeno"`     //客户订单号
	Phone       string `json:"phone"`       //电话号码
	Time        string `json:"time"`        //订购时间
	ProductId   string `json:"productid"`   //产品编号
	OrderId     string `json:"ordered"`     //系统订单编号
}

func (this *Fengzheng) QuerySingle(tradeno string) {
	request := queryRequest{}
	request.TradeNo = tradeno
	request.OrderId = 0
	response := queryResponse{}
	this.Call("query/single", &request, &response)
}

//余额查询
type queryBalanceRequest struct {
	baseParames
	Time string `json:"time"`
}

type QueryBalanceResponse struct {
	Balance string `json:"balance"`
	Status  string `json:"status"`
}

func (this *queryBalanceRequest) ToMap() map[string]string {
	result := make(map[string]string)
	result["time"] = this.Time
	return result
}

func (this *Fengzheng) QueryBalance() *QueryBalanceResponse {
	request := queryBalanceRequest{}
	request.Time = utils.Now()
	response := QueryBalanceResponse{}
	this.Call("balance", &request, &response)
	return &response
}
