package tongcheng

import (
	"github.com/aosfather/bingo_utils/http"
)

/**
获取订单详情
*/

type OrderDetailReq struct {
	Sign          string /** 参考sign加密算法[无须] **/
	DistributorCd string /** 分销商[必须] **/
	OrderNo       string //订单号
}

func (this *OrderDetailReq) toMap() map[string]string {
	values := make(map[string]string)
	values["sign"] = this.Sign
	values["distributorCd"] = this.DistributorCd
	values["orderNo"] = this.OrderNo
	return values

}

type OrderDetailResp struct {
	TcbaseResp
	Content OrderDetail `json:"content"`
}

type OrderDetail struct {
	OrderNo         string  `json:"orderNo"`
	ContactName     string  `json:"contactName"`     //联系人名称
	ContactPhone    string  `json:"contactPhone"`    //联系人手机号
	OrderStatus     string  `json:"orderStatus"`     //订单状态编号
	OrderStatusName string  `json:"orderStatusName"` //订单状态名称
	TotalAmount     float64 `json:"totalAmount"`     //订单价格
	ActualAmount    float64 `json:"actualAmount"`    //同程结算价(实际扣款金额)
	CarUseTime      int64   `json:"carUseTime"`      //
	ProductId       int64   `json:"productId"`
	StartCityId     int64   `json:"startCityId"`
	EndCityId       int64   `json:"endCityId"`
	StartAddress    string  `json:"startAddress"` //出发地址
	EndAddress      string  `json:"endAddress"`   //目的地址
	DriverName      string  `json:"driverName"`   //司机名称
	DriverPhone     string  `json:"driverTel"`    //司机手机号
	CarTypeName     string  `json:"carTypeName"`  //车型名称
	Seats           int64   `json:"seats"`        //车辆座位数
	VehicleBrand    string  `json:"vehicleBrand"` //车辆品牌
	VehicleNo       string  `json:"vehicleNo"`    //车牌号
	VehicleColor    string  `json:"vehicleColor"` //车辆颜色
	Flng            float64 `json:"startLon"`     //出发地经度
	Flat            float64 `json:"startLat"`     //出发地纬度
	Tlng            float64 `json:"endLon"`       //目的地经度
	Tlat            float64 `json:"endLat"`       //目的地纬度

	LandMark    string `json:"landmarkNo"`  //机场三字码或火车站编号
	FlightNo    string `json:"flightNo"`    //航班号
	FlightDate  int64  `json:"flightDate"`  //航班起飞日期
	Baggages    int64  `json:"baggages"`    //可放行李数
	CarTypeDesc string `json:"carTypeDesc"` //车型描述
}

func (this *Tongcheng) GetOrderDetail(orderNo string) *OrderDetailResp {
	req := OrderDetailReq{}
	req.OrderNo = orderNo
	req.DistributorCd = this.DistributorCd
	req.Sign = sign(req.toMap(), this.TcAccessToken)

	resp := OrderDetailResp{}
	err := http.Post(this.Domain+TcOrderDetailUrl, req, &resp)
	if err != nil {
		resp.Status = "FAILURE"
		resp.Code = "999"
		resp.Msg = err.Error()
	}
	return &resp
}
