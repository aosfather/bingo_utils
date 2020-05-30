package tongcheng

import (
	"fmt"
	"github.com/aosfather/bingo_utils/http"
)

type CancelOrderReq struct {
	Sign          string /** 参考sign加密算法[无须] **/
	DistributorCd string /** 分销商[必须] **/
	OrderNo       string //订单号
	ForceFlag     bool   //强制取消标志
	CancelReason  string //取消原因

}

func (this *CancelOrderReq) toMap() map[string]string {
	values := make(map[string]string)
	values["sign"] = this.Sign
	values["distributorCd"] = this.DistributorCd
	values["orderNo"] = this.OrderNo
	values["forceFlag"] = fmt.Sprintf("%s", this.ForceFlag)
	values["cancelReason"] = this.CancelReason

	return values

}

type TcbaseResp struct {
	Code   string `json:"code"`   /** 暂无 */
	Msg    string `json:"msg"`    /** 错误信息，无错误为空 */
	Status string `json:"status"` /** 状态（SUCCESS:成功，FAILURE:失败） **/
}

type CancelOrderResp struct {
	TcbaseResp
	Content CancelOrderResult `json:"content"`
}

func (this *CancelOrderResp) IsSuccess() bool {
	return this.Status == TcSuccess
}

type CancelOrderResult struct {
	PenaltyFlag bool
	CostPenalty float64
}

func (this *Tongcheng) CancelOrderSimple(orderNo string, force bool, reason string) *CancelOrderResp {
	req := CancelOrderReq{}
	req.OrderNo = orderNo
	req.ForceFlag = force
	req.CancelReason = reason
	return this.CancelOrder(&req)
}

//取消订单
func (this *Tongcheng) CancelOrder(req *CancelOrderReq) *CancelOrderResp {
	//取消订单
	req.DistributorCd = this.DistributorCd
	req.Sign = sign(req.toMap(), this.TcAccessToken)

	resp := CancelOrderResp{}
	err := http.Post(this.Domain+TcCancelUrl, req, &resp)
	if err != nil {
		resp.Status = "FAILURE"
		resp.Code = "999"
		resp.Msg = err.Error()
	}
	return &resp
}
