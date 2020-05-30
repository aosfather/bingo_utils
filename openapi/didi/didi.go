package didi

import (
	"encoding/json"
	"fmt"
	"github.com/aosfather/bingo_utils/codes"
	"github.com/aosfather/bingo_utils/http"
	"github.com/aosfather/bingo_utils/log"
	"sort"
	"time"
)

const (
	_DIDI_URL  = "http://api.es.xiaojukeji.com"
	_TOKEN_URL = _DIDI_URL + "/v1/Auth/authorize"

	ERR_SUCCESS = 0
)

var (
	_Request_LEVEL = []requestLevel{requestLevel{"600", 301}, requestLevel{"900", 301}, requestLevel{"100", 201}, requestLevel{"400", 201}, requestLevel{"200", 201}}
)

type requestLevel struct {
	level string
	rule  int
}
type DidiBaseResponse struct {
	Code    int    `json:"errno"`
	Message string `json:"errmsg"`
}
type didRequest interface {
	toMap() map[string]string
}
type didiBaseRequest struct {
	Client    string
	Token     string
	Timestamp string
	Sign      string
}

func (this *didiBaseRequest) toMap() map[string]string {
	values := make(map[string]string)
	values["client_id"] = this.Client
	values["access_token"] = this.Token
	values["timestamp"] = this.Timestamp

	return values
}

type accessTokenResponse struct {
	Token   string `json:"access_token"`
	Expires int64  `json:"expires_in"`
	Type    string `json:"token_type"`
	Scope   string `json:"scope"`
}

type DiDiConfig struct {
	accessToken string //token
	SignKey     string //
	Client      string //
	Secret      string //
	GrantType   string //
	Phone       string
	log         log.Log
	expireTime  time.Time
	Notify      NotifyProcess
}

func (this *DiDiConfig) Init() {
	this.getAccessToken()
}

func (this *DiDiConfig) callGetSimpleApi(url string, result interface{}) {
	request := didiBaseRequest{}
	this.callGetApi(url, &request, &request, result)

}

func (this *DiDiConfig) callGetApi(url string, base *didiBaseRequest, request didRequest, result interface{}) {
	if request == nil {
		return
	}
	this._checkAndRefreshToken() //检查和刷新token

	base.Client = this.Client
	base.Token = this.accessToken
	base.Timestamp = getTimestamp()
	parames := request.toMap()
	this.httpGet(url, parames, result)
}

func (this *DiDiConfig) callApi(url string, base *didiBaseRequest, request didRequest, result interface{}) {
	if request == nil {
		return
	}
	this._checkAndRefreshToken() //检查和刷新token

	base.Client = this.Client
	base.Token = this.accessToken
	base.Timestamp = getTimestamp()
	parames := request.toMap()
	sign := this.Sign(parames)
	parames["sign"] = sign
	http.Post(url, parames, result)
}

func (this *DiDiConfig) httpGet(url string, parames map[string]string, result interface{}) {
	sign := this.Sign(parames)
	parames["sign"] = sign
	theParames := ""
	for key, value := range parames {
		if theParames == "" {
			theParames = key + "=" + value
		} else {
			theParames = theParames + "&" + key + "=" + value
		}

	}

	datas, err := http.HTTPGet(url + "?" + theParames)
	this.log.Debug("http get-> %s-%s-%s", url, theParames, string(datas))
	if err == nil {
		err = json.Unmarshal(datas, result)
		if err != nil {
			this.log.Error("to json error:%s", err.Error())
		}
	} else {
		this.log.Error(err.Error())
	}

}

func (this *DiDiConfig) _checkAndRefreshToken() {
	now := time.Now()
	if now.After(this.expireTime) {
		this.getAccessToken()
	}
}

func (this *DiDiConfig) getAccessToken() {
	parames := this.toMap()
	sign := this.Sign(parames)
	parames["sign"] = sign
	result := accessTokenResponse{}
	err := http.Post(_TOKEN_URL, parames, &result)
	if err == nil {
		this.accessToken = result.Token
		this.expireTime = time.Now().Add(time.Second * time.Duration(result.Expires))
	}

}

func (this *DiDiConfig) toMap() map[string]string {
	values := make(map[string]string)
	values["client_id"] = this.Client
	values["client_secret"] = this.Secret
	values["grant_type"] = this.GrantType
	values["phone"] = this.Phone
	values["timestamp"] = getTimestamp()

	return values
}

func (this *DiDiConfig) Sign(parames map[string]string) string {
	if parames != nil {
		parames["sign_key"] = this.SignKey
	}
	nameArray := []string{}
	for key, _ := range parames {
		nameArray = append(nameArray, key)
	}

	sort.Strings(nameArray)
	//使用&来拼接字符串
	var theString = ""
	for _, name := range nameArray {
		if "sign" == name {
			continue
		}
		if theString == "" {
			theString = name + "=" + parames[name]
		} else {
			theString = theString + "&" + name + "=" + parames[name]
		}

	}
	delete(parames, "sign_key")
	return codes.ToMd5Hex(theString)
}

//------------------基本------------------------//

func getTimestamp() string {
	return fmt.Sprintf("%d", time.Now().UnixNano()/1e6)
}
