package didi

import (
	"fmt"
)

/**
didi 订单处理
1、获取订单号
2、下单
3、取消订单
4、确认支付
*/

const (
	_ORDER_ID      = _DIDI_URL + "/v1/order/Create/orderId"
	_ORDER_REQUEST = _DIDI_URL + "/v1/order/Create/request"
	_CANCEL_ORDER  = _DIDI_URL + "/v1/order/Cancel"
	_FEE_CONFIRM   = _DIDI_URL + "/v1/order/FeeConfirm"
)

//------------获取订单号---------------//

type orderResponse struct {
	didiBaseResponse
	Data orderId `json:"data"`
}

type orderId struct {
	Id string `json:"order_id"`
}

//获取滴滴生成的订单号
func (this *DiDiConfig) createOrderId() string {
	result := orderResponse{}
	this.callGetSimpleApi(_ORDER_ID, &result)
	return result.Data.Id
}

//------------创建订单-------------------//
type OrderRequest struct {
	didiBaseRequest
	OrderBase
	Order      string `json:"order_id"` //请求id 获取请参见 获取请求id
	Rule       int    //计价模型分类，201(专车)；301(快车)
	Time       string `json:"app_time"`        //客户端时间（例如：2015-06-16 12:00:09）
	DynamicMd5 string `json:"dynamic_md5"`     //价格md5,通过 新的预估价接口获得
	Lineup     int    `json:"enable_lineup"`   //是否允许排队：0,为不允许；1,为允许。默认为0
	Reassign   int    `json:"enable_reassign"` //是否允许改派：0,为不允许；1,为允许。默认为0
	/*
		enable_lineup中，订单是否会排队，由滴滴的大数据排队策略控制。该字段选择允许排队，则代表在该订单满足大数据排队策略时，自动进入队列进行排队；如不满足排队策略，即便该字段选择允许排队，订单也不会进入队列
	当enable_lineup为1（允许排队）时，enable_reassign必须传1（允许改派），否则会报错
	*/
}

func (this *OrderRequest) toMap() map[string]string {
	values := this.didiBaseRequest.toMap()
	values["order_id"] = this.Order
	values["rule"] = fmt.Sprintf("%d", this.Rule)
	values["type"] = fmt.Sprintf("%d", this.Type)
	values["passenger_phone"] = this.PassengerPhone
	//出发地
	values["city"] = this.City
	values["flat"] = this.Flat
	values["flng"] = this.Flng
	values["start_name"] = this.StartName
	values["start_address"] = this.StartAddress
	//目的地
	values["tlat"] = this.Tlat
	values["tlng"] = this.Tlng
	values["end_name"] = this.EndName
	values["end_address"] = this.EndAddress

	//时间及类型
	if this.Type == 1 {
		values["departure_time"] = this.Departure
	}
	values["require_level"] = this.Level
	values["app_time"] = this.Time
	values["dynamic_md5"] = this.DynamicMd5
	return values
}

type OrderResponse struct {
	didiBaseResponse
	Data OrderRequstResult
}

type OrderBase struct {
	Id             string //订单id
	City           string //城市id
	Type           int    //订单类型
	CallPhone      string `json:"call_phone"`
	PassengerPhone string `json:"passenger_phone"` //乘车人手机号
	Status         int    //订单状态
	Flat           string //出发地纬度
	Flng           string
	Tlat           string //目的地纬度
	Tlng           string
	Clat           string //当前纬度
	Clng           string
	StartName      string `json:"start_name"`     //出发地名称
	StartAddress   string `json:"start_address"`  //出发地地址
	EndName        string `json:"end_name"`       //目的地名称
	EndAddress     string `json:"end_address"`    //目的地地址
	Extra          string `json:"end_address"`    //目的地地址
	Departure      string `json:"departure_time"` //出发时间
	OrderTime      string `json:"order_time"`     //下单时间
	Level          string `json:"require_level"`  //所需车型
	Remark         string //备注
}
type OrderRequstResult struct {
	Order OrderBase
	Combo OrderCombo
	Price OrderPrice
}

type OrderCombo struct {
	Time     int    //套餐时长
	Distance string //套餐距离
	Fee      string //套餐价格
}

type OrderPrice struct {
	Estimate string
}

//创建订单，需要先获取订单号，然后在创建，如果创建失败返回nil，创建成功返回创建结果
func (this *DiDiConfig) CreateOrder(request *OrderRequest) *OrderResponse {
	request.Order = this.createOrderId() //请求订单号
	//检查必输参数？
	response := OrderResponse{}
	this.callApi(_ORDER_REQUEST, &request.didiBaseRequest, request, &response)
	if response.Code == ERR_SUCCESS {
		this.log.Debug("create didi order success,%v", response)

	}
	return &response

}

//----------------------------------------------------------------------//

//------------订单取消----------------//
type OrderCancelRequest struct {
	didiBaseRequest
	Order string `json:"order_id"` //订单id
	Force string //是否强制取消(true 或 false)默认false

}

func (this *OrderCancelRequest) toMap() map[string]string {
	values := this.didiBaseRequest.toMap()
	values["order_id"] = this.Order
	values["force"] = this.Force
	return values
}

type OrderCancelResponse struct {
	didiBaseResponse
	Data OrderCancelCost
}
type OrderCancelCost struct {
	Cost string
}

//取消订单，并返回需要花费的钱
func (this *DiDiConfig) CancelOrder(id string, force bool) (string, error) {
	request := OrderCancelRequest{}
	request.Order = id
	if force {
		request.Force = "true"
	} else {
		request.Force = "false"
	}

	response := OrderCancelResponse{}
	this.callApi(_CANCEL_ORDER, &request.didiBaseRequest, &request, &response)
	//判断返回状态
	if response.Code != ERR_SUCCESS {
		this.log.Debug("cancel didi order failed,%v", response)
		return "", fmt.Errorf("%d:%s", response.Code, response.Message)
	}

	return response.Data.Cost, nil

}

//-------------费用确认-------------//
type orderSimpleRequest struct {
	didiBaseRequest
	Order string `json:"order_id"` //订单id
}

func (this *orderSimpleRequest) toMap() map[string]string {
	values := this.didiBaseRequest.toMap()
	values["order_id"] = this.Order

	return values
}

func (this *DiDiConfig) FeeConfirm(id string) {
	request := orderSimpleRequest{}
	request.Order = id
	response := didiBaseResponse{}
	this.callApi(_FEE_CONFIRM, &request.didiBaseRequest, &request, &response)
	if response.Code == ERR_SUCCESS {

	}

}
