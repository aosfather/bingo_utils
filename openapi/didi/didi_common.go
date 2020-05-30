package didi

import (
	"fmt"
	"github.com/aosfather/bingo_utils/codes"
)

const (
	_ESTIMATE_URL   = _DIDI_URL + "/v1/common/Estimate/priceCoupon" //
	_GETADDRESS_URL = _DIDI_URL + "/v1/common/Address/getAddress"   //地址联想
	_ALLCITY_URL    = _DIDI_URL + "/v1/common/Cities/getAll"        //获取城市列表
	_PRICE_RULE_URL = _DIDI_URL + "/v1/common/Cities/getPrice"      //获取城市的计价规则
	_FEATURE_URL    = _DIDI_URL + "/v1/common/Estimate/getFeature"  //获取旅程时间预估
)

//获取城市列表
type DidiCity struct {
	Name     string //城市名称
	City     int    `json:"cityid"`        //城市id
	Zhuanche int    `json:"open_zhuanche"` //是否开通专车 1-开通 0-未开通
	Kuaiche  int    `json:"open_kuaiche"`  //是否开通快车 1-开通 0-未开通
}
type didiCityResponse struct {
	DidiBaseResponse
	Data []DidiCity
}

func (this *DiDiConfig) GetAllCity() []DidiCity {
	result := didiCityResponse{}
	this.callGetSimpleApi(_ALLCITY_URL, &result)
	if result.Code == ERR_SUCCESS {
		return result.Data
	}
	return nil
}

//获取地址提示
type didiAddressRequest struct {
	didiBaseRequest
	City  string
	Input string
}

func (this *didiAddressRequest) toMap() map[string]string {
	values := this.didiBaseRequest.toMap()
	values["city"] = this.City
	values["input"] = this.Input

	return values
}

type didiAddressResponse struct {
	DidiBaseResponse
	Data didiAddressData `json:"data"`
}

type didiAddressData struct {
	Input     string        `json:"input"`
	PlaceData []didiAddress `json:"place_data"`
}

type didiAddress struct {
	Name    string  `json:"displayname"`
	Address string  `json:"address"`
	City    string  `json:"city"`
	Area    int     `json:"area"`
	Lng     float32 `json:"lng"`
	Lat     float32 `json:"lat"`
}

func (this *DiDiConfig) GetAddress(request *didiAddressRequest) *didiAddressResponse {
	result := didiAddressResponse{}
	this.callGetApi(_GETADDRESS_URL, &request.didiBaseRequest, request, &result)
	return &result
}

//-----预估行程属性---//
type FeatureRequest struct {
	didiBaseRequest
	Flat string //出发地纬度
	Flng string //出发地经度
	Tlat string //目的地纬度
	Tlng string //目的地经度

	Rule int //计价模型分类，201(专车)；301(快车)
	City int //出发城市id（城市车型接口返回）
}

func (this *FeatureRequest) toMap() map[string]string {
	values := this.didiBaseRequest.toMap()
	values["flat"] = this.Flat
	values["flng"] = this.Flng
	values["tlat"] = this.Tlat
	values["tlng"] = this.Tlng
	values["city"] = fmt.Sprintf("%d", this.City)
	values["rule"] = fmt.Sprintf("%d", this.Rule)
	return values
}

type FeatureResponse struct {
	DidiBaseResponse
	Data Feature
}

type Feature struct {
	Dist     int //行驶距离 单位：米
	Duration int //行驶时间 单位：秒
	SlowTime int `json:"slow_time"` //低速时间 单位：秒
}

func (this *DiDiConfig) GetFeature(request *FeatureRequest) *Feature {
	result := FeatureResponse{}
	this.callGetApi(_FEATURE_URL, &request.didiBaseRequest, request, &result)
	if result.Code == ERR_SUCCESS {
		return &result.Data
	}

	this.log.Error("getfeature error:%d,%s", result.Code, result.Message)
	return nil
}

//------滴滴预估价-----------//
type DidiPriceRequest struct {
	didiBaseRequest
	Flat      string //出发地纬度
	Flng      string //出发地经度
	Tlat      string //目的地纬度
	Tlng      string //目的地经度
	Level     string //车型代码，详情见：订单状态定义
	Rule      int    //计价模型分类，201(专车)；301(快车)
	City      int    //出发城市id（城市车型接口返回）
	Type      int    //0:实时单 1:预约单
	Departure string //预约单必须传（格式例如：2015-06-16 12:00:09）
	Mode      int    //计价模式：0-普通计价 1-一口价 默认为0
}

func (this *DidiPriceRequest) toMap() map[string]string {
	values := this.didiBaseRequest.toMap()
	values["flat"] = this.Flat
	values["flng"] = this.Flng
	values["tlat"] = this.Tlat
	values["tlng"] = this.Tlng
	values["require_level"] = this.Level
	if this.Type == 1 {
		values["departure_time"] = this.Departure
	} else {
		//		values["departure_time"] = "2018-02-12 12:00:09"
	}
	values["city"] = fmt.Sprintf("%d", this.City)
	values["rule"] = fmt.Sprintf("%d", this.Rule)
	values["type"] = fmt.Sprintf("%d", this.Type)
	values["pricing_mode"] = fmt.Sprintf("%d", this.Mode)

	return values
}

type DidiCarPrice struct {
	Name         string  `json:"name"` //单车类型名称
	City         string  `json:"cityid"`
	Code         string  `json:"code"`              //单车类型的对应码（100舒适型，400六座商务, 200行政级, 600普通快车,900优享快车）
	Price        float32 `json:"price"`             //总价格(包含dynamic_price) 单位：元
	DynamicPrice float32 `json:"dynamic_price"`     //动调溢价 单位：元
	DynamicMd5   string  `json:"dynamic_md5"`       //价格md5
	PriceTip     string  `json:"price_tip"`         //价格提示
	StartPrice   float32 `json:"start_price"`       //起步价格 单位：元
	UnitPrice    float32 `json:"normal_unit_price"` //每公里单价 单位：元
}

type DidiPriceResponse struct {
	DidiBaseResponse
	Data  map[string]DidiCarPrice `json:"data"`
	level string
}

func (this *DiDiConfig) GetPrice(request *DidiPriceRequest) map[string]DidiCarPrice {
	theResult := make(map[string]DidiCarPrice)
	for _, v := range _Request_LEVEL {
		req := request
		req.Level = v.level
		req.Rule = v.rule
		result := DidiPriceResponse{}
		this.log.Debug("call the didi price interface")
		this.callGetApi(_ESTIMATE_URL, &req.didiBaseRequest, req, &result)
		theResult[v.level] = result.Data[v.level]

	}

	return theResult
}

//--------获取城市的计价规则--------------------//
type priceRuleRequest struct {
	didiBaseRequest
	City string
	Rule string
}

type priceRuleResponse struct {
	DidiBaseResponse
	Data map[string]*priceRuleInfo `json:"data"`
}

type priceRuleInfo struct {
	Area                      codes.AsString `json:"area"`
	District                  string         `json:"district"`
	Car_level                 string         `json:"car_level"`
	Cancel_book_money         string         `json:"cancel_book_money"`
	Cancel_book_time          string         `json:"cancel_book_time"`
	Cancel_real_money         string         `json:"cancel_real_money"`
	Cancel_real_time          string         `json:"cancel_real_time"`
	Fuel_fee                  string         `json:"fuel_fee"`
	Min_charge                string         `json:"min_charge"`
	Appointment_min_charge    codes.AsString `json:"appointment_min_charge"`
	Night_driving_unit_price  codes.AsString `json:"night_driving_unit_price"`
	Night_begin               string         `json:"night_begin"`
	Night_end                 string         `json:"night_end"`
	Empty_distance            codes.AsString `json:"empty_distance"`
	Empty_driving_unit_price  codes.AsString `json:"empty_driving_unit_price"`
	Low_speed_time_unit_price codes.AsString `json:"low_speed_time_unit_price"`
	Start_price               codes.AsString `json:"start_price"`
	Start_distance            codes.AsString `json:"start_distance"`
	Normal_unit_price         codes.AsString `json:"normal_unit_price"`
	Time_unit_price           codes.AsString `json:"time_unit_price"`
	C_level                   string         `json:"c_level"`
	C_icon                    string         `json:"c_icon"`
	C_android_icon            string         `json:"c_android_icon"`
	C_logo                    string         `json:"c_logo"`
	C_image                   string         `json:"c_image"`
	C_level_name              string         `json:"c_level_name"`
	C_level_desc              string         `json:"c_level_desc"`
}

func (this *priceRuleRequest) toMap() map[string]string {
	values := this.didiBaseRequest.toMap()
	values["city"] = this.City
	values["rule"] = this.Rule

	return values
}

func (this *DiDiConfig) GetPriceRule(city string) []*priceRuleInfo {
	var result []*priceRuleInfo
	//查询201规则
	request := priceRuleRequest{}
	request.City = city
	request.Rule = "201"
	response := priceRuleResponse{}
	this.callGetApi(_PRICE_RULE_URL, &request.didiBaseRequest, &request, &response)
	if response.Code == 0 {
		result = append(result, response.Data["100"])
		result = append(result, response.Data["200"])
		result = append(result, response.Data["400"])
	}

	//查询301规则
	request.Rule = "301"
	this.callGetApi(_PRICE_RULE_URL, &request.didiBaseRequest, &request, &response)
	if response.Code == 0 {
		result = append(result, response.Data["600"])
	}

	return result
}
