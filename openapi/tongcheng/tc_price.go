package tongcheng

import (
	"fmt"
	"github.com/aosfather/bingo_utils/http"
	"time"

	"xingyun.com/framework"
)

type EstimatePriceListResp struct {
	TcResp
	Content []EstimatePriceResp `json:"content"`
}

type EstimatePriceReq struct {
	Sign               string /** 参考sign加密算法[无须] **/
	DistributorCd      string /** 分销商[无须] **/
	ProductId          string /** 产品ID[必须]，参考产品列表附录 **/
	StartCityId        string /** 出发城市id[必须] **/
	EndCityId          string /** 目的地城市id[必须] **/
	StartAddress       string /** 出发地地址[必须] **/
	StartAddressDetail string /** 出发地详细地址[必须] **/
	StartLongitude     string /** 出发地经度[必须] **/
	StartLatitude      string /** 出发地纬度[必须] **/
	EndAddress         string /** 目的地详细地址[必须] **/
	EndAddressDetail   string /** 目的地详细地址[必须] **/
	EndLongitude       string /** 目的地经度[必须] **/
	EndLatitude        string /** 目的地纬度[必须] **/
	CarUseTime         string /** 用车时间(unix时间戳)[必须] **/
	LandmarkNo         string /** 机场三字码或火车站编号 **/
	FlightNo           string /** 航班号 **/
	FlightDate         string /** 航班起飞日期 **/
}

func (this *EstimatePriceReq) toMap() map[string]string {
	values := make(map[string]string)
	values["sign"] = this.Sign
	values["distributorCd"] = this.DistributorCd
	values["productId"] = this.ProductId
	values["startCityId"] = this.StartCityId
	values["endCityId"] = this.EndCityId
	values["startAddress"] = this.StartAddress
	values["startAddressDetail"] = this.StartAddressDetail
	values["startLongitude"] = this.StartLongitude
	values["startLatitude"] = this.StartLatitude
	values["endAddress"] = this.EndAddress
	values["endAddressDetail"] = this.EndAddressDetail
	values["endLongitude"] = this.EndLongitude
	values["endLatitude"] = this.EndLatitude
	values["carUseTime"] = this.CarUseTime
	values["landmarkNo"] = this.LandmarkNo
	values["flightNo"] = this.FlightNo
	values["flightDate"] = this.FlightDate

	return values

}

func (this *EstimatePriceReq) vaild() BaseResp {
	var result = BaseResp{}
	if nil == this {
		fmt.Println("method:validEstimatePriceReq,error:请求对象不能为空")
		result.Error()
		result.Message = "请求对象不能为空"
		return result
	}
	if "" == this.ProductId {
		fmt.Println("method:validEstimatePriceReq,error:产品ID不能为空")
		result.Error()
		result.Message = "产品ID不能为空"
		return result
	}

	if "" == this.StartCityId {
		fmt.Println("method:validEstimatePriceReq,error:出发城市id不能为空")
		result.Error()
		result.Message = "出发城市id不能为空"
		return result
	}
	if "" == this.EndCityId {
		fmt.Println("method:validEstimatePriceReq,error:目的地城市id不能为空")
		result.Error()
		result.Message = "目的地城市id不能为空"
		return result
	}
	if "" == this.StartAddress {
		fmt.Println("method:validEstimatePriceReq,error:出发地地址不能为空")
		result.Error()
		result.Message = "出发地地址不能为空"
		return result
	}
	if "" == this.StartAddressDetail {
		fmt.Println("method:validEstimatePriceReq,error:出发地地址详情不能为空")
		result.Error()
		result.Message = "出发地地址详情不能为空"
		return result
	}
	if "" == this.StartLongitude {
		fmt.Println("method:validEstimatePriceReq,error:出发地经度不能为空")
		result.Error()
		result.Message = "出发地经度不能为空"
		return result
	}
	if "" == this.StartLatitude {
		fmt.Println("method:validEstimatePriceReq,error:出发地纬度不能为空")
		result.Error()
		result.Message = "出发地纬度不能为空"
		return result
	}
	if "" == this.EndAddress {
		fmt.Println("method:validEstimatePriceReq,error:目的地地址不能为空")
		result.Error()
		result.Message = "目的地地址不能为空"
		return result
	}
	if "" == this.EndAddressDetail {
		fmt.Println("method:validEstimatePriceReq,error:目的地地址详情不能为空")
		result.Error()
		result.Message = "目的地地址详情不能为空"
		return result
	}
	if "" == this.EndLongitude {
		fmt.Println("method:validEstimatePriceReq,error:目的地经度不能为空")
		result.Error()
		result.Message = "目的地经度不能为空"
		return result
	}
	if "" == this.EndLatitude {
		fmt.Println("method:validEstimatePriceReq,error:目的地纬度不能为空")
		result.Error()
		result.Message = "目的地纬度不能为空"
		return result
	}
	if "" == this.CarUseTime {
		fmt.Println("method:validEstimatePriceReq,error:用车时间不能为空")
		result.Error()
		result.Message = "用车时间不能为空"
		return result
	}

	result.Success()
	return result
}

type EstimatePriceResp struct {
	SupplierCd        string  /** 供应商编号 **/
	SupplierName      string  /** 供应商名称 **/
	CarTypeId         int     /** 同程车型ID **/
	CarTypeName       string  /** 同程车型名称 **/
	Seats             int     /** 车辆座位数 **/
	Price             float32 /** 预估价格 **/
	Baggages          int     /** 可放行李数 **/
	CarTypeDesc       string  /** 车型描述 **/
	PriceEstimateMark string  /** 价格预估标志，用来标明本次约定的价格（下同程单时使用，30分钟有效期） **/
}

func (this *Tongcheng) EstimatePriceQuery(estimatePriceReq EstimatePriceReq) (BaseResp, []EstimatePriceResp) {

	//#1.校验请求参数
	var baseResp = estimatePriceReq.vaild()

	if !baseResp.IsSuccess() {
		//参数校验不通过,
		return baseResp, nil
	}
	//#2
	estimatePriceReq.DistributorCd = this.DistributorCd
	estimatePriceReq.Sign = sign(estimatePriceReq.toMap(), this.TcAccessToken)

	//#3.
	estimatePriceListResp := EstimatePriceListResp{}
	err := http.Post(this.Domain+TcEstimateUrl, estimatePriceReq, &estimatePriceListResp)
	if err != nil {
		this.logger.Debug(err.Error())
		baseResp.Error()
		baseResp.Message = "网络请求异常"
		return baseResp, nil
	}
	fmt.Println(estimatePriceListResp.TcResp.Status)
	if !estimatePriceListResp.TcResp.isSuccess() {
		baseResp.Error()
		baseResp.Message = estimatePriceListResp.TcResp.Msg
		return baseResp, nil
	}

	//#4
	baseResp.Success()

	return baseResp, estimatePriceListResp.Content
}

func (this *Tongcheng) EstimatePriceQueryMerge(estimatePriceReqList []EstimatePriceReq, resultChans []chan EstimatePriceListResp) EstimatePriceListResp {
	for i, estimatePriceReq := range estimatePriceReqList {
		go func(request EstimatePriceReq, resultChans chan EstimatePriceListResp) {
			_baseResp, _epList := this.EstimatePriceQuery(request)
			if !_baseResp.IsSuccess() {
				this.logger.Debug("同程预估价错误:%s", _baseResp.Message)
				return
			}
			_estimatePriceListResp := EstimatePriceListResp{}
			_estimatePriceListResp.Content = _epList
			resultChans <- _estimatePriceListResp
		}(estimatePriceReq, resultChans[i])
	}

	estimatePriceListResp := EstimatePriceListResp{}
	for _, resultChan := range resultChans {
		select {
		case result := <-resultChan:
			for _, _content := range result.Content {
				estimatePriceListResp.Content = append(estimatePriceListResp.Content, _content)
			}
			close(resultChan)
		case <-time.After(5 * time.Second):
			this.logger.Debug("预估价读取超时")
			close(resultChan)
		}
	}

	estimatePriceListResp.HasSuccess()
	return estimatePriceListResp

}
