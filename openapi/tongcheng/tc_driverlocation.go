package tongcheng

import (
	"github.com/aosfather/bingo_utils/http"
)

/**
获取司机位置
*/

type driverLocationRequest struct {
	OrderNo string `json:"tcOrderNo"`
}

type tcCommonResponse struct {
	Code    string `json:"code"`
	Message string `json:"msg"`
	Success bool   `json:"isSuccess"`
}
type DriverLocationResponse struct {
	tcCommonResponse
	Result DriverLocation `json:"result"`
}

//司机及司机定位定义信息
type DriverLocation struct {
	OrderStatus    bool    `json:"orderStatus"`    //订单状态
	Lng            float64 `json:"longitude"`      //司机经度
	Lat            float64 `json:"latitude"`       //司机纬度
	Name           string  `json:"driverName"`     //司机名字
	VehicleLicense string  `json:"vehicleLicense"` //司机车牌
	Phone          string  `json:"driverPhone"`    //司机手机号
	VehicleModel   string  `json:"vehicleModel"`   //车型
	Rate           string  `json:"driverRate"`     //司机评分
	Photo          string  `json:"driverPhoto"`    //司机头像
}

//获取司机的位置
func (this *Tongcheng) GetDriverLocation(orderno string) *DriverLocationResponse {
	resp := DriverLocationResponse{}
	if orderno == "" {
		resp.Success = false
		resp.Code = "999"
		resp.Message = "please input order no!"
		return &resp
	}

	err := http.Post(this.Domain+TcDriverLocationUrl, driverLocationRequest{orderno}, &resp)
	if err != nil {
		resp.Success = false
		resp.Code = "999"
		resp.Message = err.Error()
	}

	return &resp

}
