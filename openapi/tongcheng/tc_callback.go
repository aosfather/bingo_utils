package tongcheng

const (
	/*
		1-确认收单
		2-下单成功
		3-司机接单
		4-接单失败
		5-司机到达
		6-行程中
		7-行程结束
		8-结算价已确定
		9-订单取消
		10-异常单
		11-司机改派
		13-确认接单
		14-下单失败
	*/
	OS_APPLY          = "1"
	OS_ORDER_SUCESS   = "2"
	OS_DRIVER_APPLY   = "3"
	OS_FAILED         = "4"
	OS_DRIVER_ARRIVED = "5"
	OS_PROGRESS       = "6"
	OS_FINISH         = "7"
	OS_PRICE_CLEAR    = "8"
	OS_CANCELED       = "9"
	OS_ERROR          = "10"
	OS_DRIVER_CHANGED = "11"
	OS_PLAY           = "13" //确认接单
	OS_ORDER_FAILED   = "14" //

)

type TongChengDriverInfo struct {
	DriverName   string `json:"driverName"`   //司机名
	DriverTel    string `json:"driverTel"`    //司机电话
	VehicleBrand string `json:"vehicleBrand"` //车型
	VehicleNo    string `json:"vehicleNo"`    // 车牌号
	VehicleColor string `json:"vehicleColor"` //颜色
}

type TongChengNoticeRequestDto struct {
	Sign        string              `json:"sign"`    //签名
	OrderNo     string              `json:"orderNo"` //同程订单号
	OrderStatus string              `json:"orderStatus"`
	Driver      TongChengDriverInfo `json:"driver"`
}

func (this *TongChengNoticeRequestDto) toMap() map[string]string {
	values := make(map[string]string)

	return values
}

func (this *Tongcheng) Callback(p *TongChengNoticeRequestDto) *TcbaseResp {
	msg := ""
	if p != nil {
		targetSign := p.Sign //sign(p.toMap(),this.TcAccessToken)
		if targetSign == p.Sign {
			//通知处理
			if this.Notify(p) {
				return &TcbaseResp{"0", "success!", TcSuccess}
			} else { //通知处理失败
				msg = "notify process error!"
			}

		} else { //通知校验失败
			msg = "validate sign failed!"
		}
	}

	return &TcbaseResp{"100", msg, TcError}
}
