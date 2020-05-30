package openapi

import (
	"encoding/json"
	"fmt"
	"github.com/aosfather/bingo_utils/http"
	"strings"
)

const (
	SUCCESS_CODE = "200"
	ERROR_CODE   = "-1"
	_GAO_DE_URL  = "http://restapi.amap.com/v3/assistant/inputtips?output=json&city=%s&citylimit=true&keywords=%s&key=%s"
	_BAI_DU_URL  = "http://api.map.baidu.com/place/search?&query=%s&region=%s&output=json&key=%s"
)

var (
	//高德的key，因为高德的免费key有次数限制所以需要根据情况申请多个轮流来使用
	_GAO_DE_KEY []string
	_Baidu_Key  = "bqApldE1oh6oBb98VYyIfy9S"

	_gaodekey_index = 0
)

func SetBaiduKey(key string) {
	_Baidu_Key = key
}

func SetGaodekey(key ...string) {
	_GAO_DE_KEY = key
}

type gaodeResponse struct {
	Status  string         `json:"status"`
	Count   string         `json:"count"`
	Message string         `json:"info"`
	Code    string         `json:"infocode"`
	Tips    []gaodeAddress `json:"tips"`
}
type baiduResponse struct {
	Status string         `json:"status"`
	Result []baiduAddress `json:"results"`
}
type gaodeGeoLocation GeoLocation

type GeoLocation struct {
	Lat string //纬度 `json:"lat,float32"`
	Lng string //经度 `json:"lng,float32"`
}

func (this *GeoLocation) UnmarshalJSON(data []byte) error {
	mapps := make(map[string]float32)
	json.Unmarshal(data, &mapps)
	this.Lat = fmt.Sprintf("%f", mapps["lat"])
	this.Lng = fmt.Sprintf("%f", mapps["lng"])
	return nil
}

type gaodeAddress struct {
	Id       string           `json:"id,omitempty"`
	Name     string           `json:"name"`
	District string           `json:"district"`
	Code     string           `json:"adcode"`
	Location gaodeGeoLocation `json:"location,omitempty"`
	Address  string           `json:"address,omitempty"`
	Type     string           `json:"typecode"`
}

type baiduAddress struct {
	Id        string      `json:"uid"`
	Name      string      `json:"name"`
	Tag       string      `json:"tag"`
	Url       string      `json:"detail_url"`
	Location  GeoLocation `json:"location,omitempty"`
	Address   string      `json:"address,omitempty"`
	Telephone string      `json:"telephone"`
}

func (this *gaodeGeoLocation) UnmarshalJSON(data []byte) error {
	str := string(data)

	if len(str) > 2 {
		str = str[1:]
		str = str[:len(str)-1]
		locs := strings.Split(str, ",")
		if len(locs) == 2 {

			this.Lng = locs[0]
			this.Lat = locs[1]
		}
	}

	return nil
}

func queryFromGAODE(city, input string) *gaodeResponse {
	if len(_GAO_DE_KEY) == 0 {
		return nil
	}
	_gaodekey_index++
	if _gaodekey_index > 99999999 {
		_gaodekey_index = 0
	}

	keyIndex := _gaodekey_index % len(_GAO_DE_KEY)
	result := gaodeResponse{}
	http.HTTPGetStruct(&result, _GAO_DE_URL, city, input, _GAO_DE_KEY[keyIndex])
	return &result
}

func queryFromBaidu(city, input string) *baiduResponse {
	result := baiduResponse{}
	http.HTTPGetStruct(&result, _BAI_DU_URL, input, city, _Baidu_Key)
	return &result
}

type LocationQueryResponse struct {
	Status  string           `json:"status"`
	Message string           `json:"message"`
	Result  []locationResult `json:"result"`
}

type locationResult struct {
	Input    string `json:"keyWord"`
	Name     string `json:"displayName"`
	Address  string `json:"address"`
	CityCode string `json:"cityCode"`
	Lat      string `json:"adressLat"`
	Lng      string `json:"adressLng"`
}

func QueryLocationTips(city, citycode, input string) *LocationQueryResponse {
	gaodeResult := queryFromGAODE(city, input)
	response := LocationQueryResponse{}
	response.Status = SUCCESS_CODE

	if gaodeResult == nil || len(gaodeResult.Tips) <= 0 {
		//查询百度
		baiduResult := queryFromBaidu(city, input)
		//转换百度的结果
		for _, item := range baiduResult.Result {
			if item.Location.Lat == "" || item.Location.Lng == "" {
				continue
			}
			l := locationResult{}
			l.Input = input
			l.CityCode = citycode
			l.Name = item.Name
			l.Address = item.Address
			l.Lat = item.Location.Lat
			l.Lng = item.Location.Lng
			response.Result = append(response.Result, l)
		}

	} else {
		//转换高德的结果
		for _, item := range gaodeResult.Tips {
			if item.Location.Lat == "" || item.Location.Lng == "" {
				continue
			}
			l := locationResult{}
			l.Input = input
			l.CityCode = citycode
			l.Name = item.Name
			l.Address = item.Address
			l.Lat = item.Location.Lat
			l.Lng = item.Location.Lng
			response.Result = append(response.Result, l)
		}
	}

	if response.Result == nil {
		response.Status = ERROR_CODE
		response.Message = "not found!"
	}

	return &response

}
