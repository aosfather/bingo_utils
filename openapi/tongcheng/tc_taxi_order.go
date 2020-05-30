package tongcheng

import (
	"fmt"
	"github.com/aosfather/bingo_utils/http"

	"xingyun.com/framework"
)

type CreateTaxiCarOrderResp struct {
	TcResp
	OrderNo string `json:"orderNo"`
}

type CreateTaxiCarOrderReq struct {
	Sign               string /** 参考sign加密算法[无须] **/
	DistributorCd      string /** 分销商[必须] **/
	DistributorNo      string /** 分销商流水号[无须] **/
	ContactName        string /**  联系人姓名[必须] **/
	ContactPhone       string /** 联系人手机[必须] **/
	ProductId          string /** 产品id[必须] **/
	UseTime            int64  /** 用车时间[必须] **/
	StartCityId        string /** 出发城市id[必须] **/
	StartAddress       string /** 出发地地址[必须] **/
	StartAddressDetail string /** 出发地详细地址[必须] **/
	StartLongitude     string /** 出发地经度[必须] **/
	StartLatitude      string /** 出发地纬度[必须] **/
	EndCityId          string /** 目的地城市id[必须] **/
	EndAddress         string /** 目的地地址[必须] **/
	EndAddressDetail   string /** 目的地详细地址[必须] **/
	EndLongitude       string /** 目的地经度[必须] **/
	EndLatitude        string /** 目的地纬度[必须] **/
}

func (this *CreateTaxiCarOrderReq) toMap() map[string]string {
	values := make(map[string]string)
	values["sign"] = this.Sign
	values["distributorCd"] = this.DistributorCd
	values["distributorNo"] = this.DistributorNo
	values["contactName"] = this.ContactName
	values["contactPhone"] = this.ContactPhone
	values["productId"] = this.ProductId
	values["useTime"] = fmt.Sprintf("%d", this.UseTime)
	values["startCityId"] = this.StartCityId
	values["startAddress"] = this.StartAddress
	values["startAddressDetail"] = this.StartAddressDetail
	values["startLongitude"] = this.StartLongitude
	values["startLatitude"] = this.StartLatitude
	values["endCityId"] = this.EndCityId
	values["endAddress"] = this.EndAddress
	values["endAddressDetail"] = this.EndAddressDetail
	values["endLongitude"] = this.EndLongitude
	values["endLatitude"] = this.EndLatitude

	return values

}

func (this *CreateTaxiCarOrderReq) vaild() BaseResp {
	var result = BaseResp{}
	if nil == this {
		fmt.Println("method:validCreateTaxiCarOrderReq,error:请求对象不能为空")
		result.Error()
		result.Message = "请求对象不能为空"
		return result
	}
	if "" == this.ContactName {
		fmt.Println("method:validCreateTaxiCarOrderReq,error:联系人姓名不能为空")
		result.Error()
		result.Message = "联系人姓名不能为空"
		return result
	}
	if "" == this.ContactPhone {
		fmt.Println("method:validCreateTaxiCarOrderReq,error:联系人手机不能为空")
		result.Error()
		result.Message = "联系人手机不能为空"
		return result
	}

	if "" == this.ProductId {
		fmt.Println("method:validCreateTaxiCarOrderReq,error:产品id不能为空")
		result.Error()
		result.Message = "产品id不能为空"
		return result
	}

	if 0 == this.UseTime {
		fmt.Println("method:validCreateTaxiCarOrderReq,error:用车时间不能为空")
		result.Error()
		result.Message = "用车时间不能为空"
		return result
	}
	if "" == this.StartCityId {
		fmt.Println("method:validCreateTaxiCarOrderReq,error:出发城市id不能为空")
		result.Error()
		result.Message = "出发城市id不能为空"
		return result
	}
	if "" == this.StartAddress {
		fmt.Println("method:validCreateTaxiCarOrderReq,error:出发地地址不能为空")
		result.Error()
		result.Message = "出发地地址不能为空"
		return result
	}

	if "" == this.StartAddressDetail {
		fmt.Println("method:validCreateTaxiCarOrderReq,error:出发地详细地址不能为空")
		result.Error()
		result.Message = "出发地详细地址不能为空"
		return result
	}
	if "" == this.StartLongitude {
		fmt.Println("method:validCreateTaxiCarOrderReq,error:出发地经度不能为空")
		result.Error()
		result.Message = "出发地经度不能为空"
		return result
	}

	if "" == this.StartLatitude {
		fmt.Println("method:validCreateTaxiCarOrderReq,error:出发地纬度不能为空")
		result.Error()
		result.Message = "出发地纬度不能为空"
		return result
	}
	if "" == this.EndCityId {
		fmt.Println("method:validCreateTaxiCarOrderReq,error:目的地城市id不能为空")
		result.Error()
		result.Message = "目的地城市id不能为空"
		return result
	}
	if "" == this.EndAddress {
		fmt.Println("method:validCreateTaxiCarOrderReq,error:目的地地址不能为空")
		result.Error()
		result.Message = "目的地地址不能为空"
		return result
	}

	if "" == this.EndAddressDetail {
		fmt.Println("method:validCreateTaxiCarOrderReq,error:目的地详细地址不能为空")
		result.Error()
		result.Message = "目的地详细地址不能为空"
		return result
	}
	if "" == this.EndLongitude {
		fmt.Println("method:validCreateTaxiCarOrderReq,error:目的地经度不能为空")
		result.Error()
		result.Message = "目的地经度不能为空"
		return result
	}
	if "" == this.EndLatitude {
		fmt.Println("method:validCreateTaxiCarOrderReq,error:目的地纬度不能为空")
		result.Error()
		result.Message = "目的地纬度不能为空"
		return result
	}

	result.Success()
	return result
}

func (this *Tongcheng) CreateTaxiCarOrder(createTaxiCarOrderReq CreateTaxiCarOrderReq) (BaseResp, CreateTaxiCarOrderResp) {

	createTaxiCarOrderResp := CreateTaxiCarOrderResp{}
	//#1.校验请求参数
	var baseResp = createTaxiCarOrderReq.vaild()

	if !baseResp.IsSuccess() {
		//参数校验不通过,
		return baseResp, createTaxiCarOrderResp
	}

	//#2
	createTaxiCarOrderReq.DistributorCd = this.DistributorCd
	createTaxiCarOrderReq.Sign = sign(createTaxiCarOrderReq.toMap(), this.TcAccessToken)

	//#3.
	err := http.Post(this.Domain+TcTaxiOrderCreateUrl, createTaxiCarOrderReq, &createTaxiCarOrderResp)
	if err != nil {
		baseResp.Error()
		baseResp.Message = "网络请求异常"
		return baseResp, createTaxiCarOrderResp
	}

	if !createTaxiCarOrderResp.TcResp.isSuccess() {
		baseResp.Error()
		baseResp.Message = createTaxiCarOrderResp.TcResp.Msg
		return baseResp, createTaxiCarOrderResp
	}

	//#4
	baseResp.Success()

	return baseResp, createTaxiCarOrderResp
}
