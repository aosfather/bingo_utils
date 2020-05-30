package tongcheng

import (
	"fmt"
	"github.com/aosfather/bingo_utils/http"
)

type CreatePriCarOrderResp struct {
	TcResp
	OrderNo string `json:"orderNo"`
}

type CreatePriCarOrderReq struct {
	Sign              string /** 参考sign加密算法[无须] **/
	DistributorCd     string /** 分销商[必须] **/
	PriceEstimateMark string /** 价格预估标志[必须] **/
	DistributorNo     string /** 分销商流水号[无须] **/
	FlightNo          string /** 航班号 **/
	ContactName       string /**  联系人姓名[必须] **/
	ContactPhone      string /** 联系人手机[必须] **/
}

func (this *CreatePriCarOrderReq) toMap() map[string]string {
	values := make(map[string]string)
	values["sign"] = this.Sign
	values["distributorCd"] = this.DistributorCd
	values["priceEstimateMark"] = this.PriceEstimateMark
	values["distributorNo"] = this.DistributorNo
	values["flightNo"] = this.FlightNo
	values["contactName"] = this.ContactName
	values["contactPhone"] = this.ContactPhone

	return values

}

func (this *CreatePriCarOrderReq) vaild() BaseResp {
	var result = BaseResp{}
	if nil == this {
		fmt.Println("method:validCreatePriCarOrderReq,error:请求对象不能为空")
		result.Error()
		result.Message = "请求对象不能为空"
		return result
	}
	if "" == this.PriceEstimateMark {
		fmt.Println("method:validCreatePriCarOrderReq,error:价格预估标志不能为空")
		result.Error()
		result.Message = "价格预估标志不能为空"
		return result
	}
	if "" == this.ContactName {
		fmt.Println("method:validCreatePriCarOrderReq,error:联系人姓名不能为空")
		result.Error()
		result.Message = "联系人姓名不能为空"
		return result
	}
	if "" == this.ContactPhone {
		fmt.Println("method:validCreatePriCarOrderReq,error:联系人手机不能为空")
		result.Error()
		result.Message = "联系人手机不能为空"
		return result
	}

	result.Success()
	return result
}

func (this *Tongcheng) CreatePriCarOrder(createPriCarOrderReq CreatePriCarOrderReq) (BaseResp, CreatePriCarOrderResp) {

	createPriCarOrderResp := CreatePriCarOrderResp{}
	//#1.校验请求参数
	var baseResp = createPriCarOrderReq.vaild()

	if !baseResp.IsSuccess() {
		//参数校验不通过,
		return baseResp, createPriCarOrderResp
	}

	//#2
	createPriCarOrderReq.DistributorCd = this.DistributorCd
	createPriCarOrderReq.Sign = sign(createPriCarOrderReq.toMap(), this.TcAccessToken)

	//#3.
	err := http.Post(this.Domain+TcPriOrderCreateUrl, createPriCarOrderReq, &createPriCarOrderResp)
	if err != nil {
		baseResp.Error()
		baseResp.Message = "网络请求异常"
		return baseResp, createPriCarOrderResp
	}

	if !createPriCarOrderResp.TcResp.isSuccess() {
		baseResp.Error()
		baseResp.Message = createPriCarOrderResp.TcResp.Msg
		return baseResp, createPriCarOrderResp
	}

	//#4
	baseResp.Success()

	return baseResp, createPriCarOrderResp
}
