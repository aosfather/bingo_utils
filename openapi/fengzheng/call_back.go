package fengzheng

/**
交易信息回调
*/

type callbackRequest struct {
	Id      string `json:"orderid"`
	TradeNo string `json:"tradeno"`
	Status  string `json:"status"`
	Sign    string `json:"sign"`
}

func (this *callbackRequest) ToMap() map[string]string {
	result := make(map[string]string)
	result["orderid"] = this.Id
	result["tradeno"] = this.TradeNo
	result["status"] = this.Status
	return result
}

func (this *Fengzheng) CallBack(r *callbackRequest, signSuccess func(), signFail func()) {
	if r != nil {
		tsign := this.sign(r.ToMap())
		if r.Sign == tsign {
			//成功后通知状态
			if signSuccess != nil {
				signSuccess()
			}
		}

	} else { //验证签名失败
		if signFail != nil {
			signFail()
		}

	}
}
