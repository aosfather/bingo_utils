package didi

/**
  滴滴订单详情

*/

const (
	_ORDER_DETAIL = _DIDI_URL + "/v1/order/Detail/getOrderDetail"
	SV_DIDI       = "didi"
	/** 订单状态 status
	300 	等待应答
	311 	订单超时
	400 	等待接驾
	410 	司机已到达
	500 	行程中
	600 	行程结束
	610 	行程异常结束
	700 	已支付
	*/
	OS_WAITING        = "300"
	OS_OUTTIME        = "311"
	OS_WAIT_DRIVER    = "400"
	OS_DRIVER_ARRIVED = "410"
	OS_DOING          = "500"
	OS_FINISH         = "600"
	OS_OVER           = "610"
	OS_PAYED          = "700"

	/** 子状态
		状态码（sub_status) 	描述（sub_status_tips)
	0 	未知
	3000 	等待抢单
	3001 	有司机抢单
	3002 	确定某个司机抢单，但是需要进入协商
	3003 	协商状态
	3101 	订单超时
	4000 	改派中
	4001 	等待接驾
	4002 	司机迟到
	4004 	乘客迟到
	4005 	迟到计费
	4101 	司机到达
	5000 	服务中/计费中
	6001 	正常订单待支付
	6002 	取消行程待支付
	6100 	取消订单产生费用待支付
	6101 	取消订单无需支付
	6102 	取消订单已支付
	6103 	客服关闭
	6104 	改派关闭
	6105 	未能完成服务关闭
	7000 	已完成
	*/
)

type OrderDetailResponse struct {
	DidiBaseResponse
	Data OrderDetails
}
type OrderDetails struct {
	Order    OrderDetail
	IsLineup int          `json:"is_lineup"` //当前状态是否排队，(0:否，1:是)
	Lineup   lineupInfo   `json:"lineup_info"`
	Reassign reassignInfo `json:"reassign_info"`
	Price    costDetail   `json:"price"`
}

type OrderDetail struct {
	OrderBase
	driverInfo
	realInfo
	SubStatus string //子状态

}

//司机信息
type driverInfo struct {
	DriverName       string `json:"driver_name"`
	DriverPhone      string `json:"driver_phone"`
	DriverPhoneReal  string `json:"driver_phone_real"`
	DriverNum        int    `json:"driver_num"`
	DriverCar        string `json:"driver_car_type"`
	DriverCarColor   string `json:"driver_car_color"`
	DriverCard       string `json:"driver_card"`
	DriverAvatar     string `json:"driver_avatar"`
	DriverOrderCount string `json:"driver_order_count"`
	DriverLevel      string `json:"driver_level"`
}

//改派信息
type reassignInfo struct {
	PreOrder    string `json:"pre_order_id"`    //改派前订单id，即此订单由哪个订单id改派而生成（当值为0时，表示该订单不是因为改派而生成的）
	NextOrder   string `json:"next_order_id"`   //改派后订单id，即由于订单改派而产生的订单id（当值为0时，表示当前订单未被改派而产生新订单）
	InitOrder   string `json:"init_order_id"`   //第一个被改派的订单id
	LatestOrder string `json:"latest_order_id"` //最新被指派的订单id
}

//排队信息
type lineupInfo struct {
	Ranking int `json:"ranking"`      //当前在队列里的位置
	Length  int `json:"queue_length"` //	队列总长度
	Wait    int `json:"wait_time"`    //预估需要等待的时间(秒)
}

//实时信息
type realInfo struct {
	Dlng            string //司机当前实时经度
	Dlat            string //司机当前实时维度
	OrderTime       string `json:"order_time"`        //下单时间
	StriveTime      string `json:"strive_time"`       //司机接单时间
	BeginChargeTime string `json:"begin_charge_time"` //开始计价时间
	FinishTime      string `json:"finish_time"`       //行程结束时间
	DelayTime       string `json:"delay_time_start"`  //迟到计费时间
	NormalDistance  string `json:"normal_distance"`   //实际行驶公里数
	NormalTime      string `json:"normal_time"`       //实际行驶时长（分钟)
	StriveLevel     string `json:"strive_level"`      //叫单车型（100舒适型，400六座商务, 200行政级,600普通快车,900优享快车）
}
type costDetail struct {
	Total string `json:"total_price"`

	Details []costInfo `json:"detail"`
}

//费用
type costInfo struct {
	Name   string //费用名称
	Amount string //金额
	Type   string //类型
}

func (this *DiDiConfig) GetOrderDetail(id string) *OrderDetailResponse {
	request := orderSimpleRequest{}
	request.Order = id
	response := OrderDetailResponse{}
	this.callGetApi(_ORDER_DETAIL, &request.didiBaseRequest, &request, &response)
	return &response
}

//司机位置信息
type DriverLocation struct {
	OrderStatus    int    //订单状态
	Lng            string //司机经度
	Lat            string //司机纬度
	Name           string //司机名字
	VehicleLicense string //司机车牌
	Phone          string //司机手机号
	VehicleModel   string //车型
	Rate           string //司机评分
	Photo          string //司机头像
}

//获取滴滴司机的位置信息
func (this *DiDiConfig) GetDriverLocation(id string) *DriverLocation {
	//查询订单详情，获取司机的信息及实时位置
	response := this.GetOrderDetail(id)
	if response.Code == ERR_SUCCESS {
		location := DriverLocation{}
		order := response.Data.Order
		//订单状态
		location.OrderStatus = order.Status
		//获取司机信息
		location.Name = order.DriverName
		location.Phone = order.DriverPhone
		location.VehicleLicense = order.DriverCard
		location.VehicleModel = order.DriverCar
		location.Rate = order.DriverLevel
		location.Photo = order.DriverAvatar
		//获取司机坐标信息
		location.Lng = order.Dlng
		location.Lng = order.Dlat

		return &location
	}

	return nil
}
