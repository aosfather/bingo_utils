package duiba

import (
	"github.com/aosfather/bingo_utils/codes"
	"sort"
)

/**

对接兑吧的接口
1、兑吧扣款回调
2、兑吧兑换成功和失败的通知
*/

const (
	DUIBA_SUCCESS = "ok"
	DUIBA_FAIL    = "fail"
)

type DuibaResponse struct {
	Status  string `json:"status"`       //取值ok 或 fail
	Message string `json:"errorMessage"` //消息
	Id      string `json:"bizId"`        //业务id
	Credits int64  `json:"credits"`      //余额
}

//兑吧扣款请求
type DuibaExchangeRequst struct {
	Uid         string `Field:"uid"`         //用户id
	Amount      int64  `Field:"credits"`     //本次兑换扣除的积分
	ItemCode    string `Field:"itemCode"`    //自有商品商品编码(非必须字段)
	App         string `Field:"appKey"`      //接口appKey，应用的唯一标识
	Date        string `Field:"timestamp"`   //1970-01-01开始的时间戳，毫秒为单位。
	Description string `Field:"description"` //本次积分消耗的描述(带中文，请用utf-8进行url解码)
	OrderNum    string `Field:"orderNum"`    //兑吧订单号(请记录到数据库中)
	Type        string `Field:"type"`        //兑换类型：alipay(支付宝), qb(Q币), coupon(优惠券), object(实物), phonebill(话费), phoneflow(流量), virtual(虚拟商品), littleGame（小游戏）,singleLottery(单品抽奖)，hdtoolLottery(活动抽奖),htool(新活动抽奖),manualLottery(手动开奖),ngameLottery(新游戏),questionLottery(答题),quizzLottery(测试题),guessLottery(竞猜) 所有类型不区分大小写
	FacePrice   string `Field:"facePrice"`   //兑换商品的市场价值，单位是分，请自行转换单位
	ActualPrice string `Field:"actualPrice"` //此次兑换实际扣除开发者账户费用，单位为分
	Ip          string `Field:"ip"`          //用户ip，不保证获取到

	Audit  string `Field:"waitAudit"` //是否需要审核(如需在自身系统进行审核处理，请记录下此信息)
	Params string `Field:"params"`    //详情参数，不同的类型，返回不同的内容，中间用英文冒号分隔
	Sign   string `Field:"sign"`      //MD5签名，详见签名规则

}

func (this *DuibaExchangeRequst) ToMap() map[string]string {
	result := make(map[string]string)
	result["uid"] = this.Uid
	return result
}

//2、兑吧兑换结果通知
type DuibaExchangeNotify struct {
	App      string `Field:"appKey"`      //接口appKey，应用的唯一标识
	Date     string `Field:"timestamp"`   //1970-01-01开始的时间戳，毫秒为单位。
	Uid      string `Field:"uid"`         //用户id
	OrderNum string `Field:"orderNum"`    //兑吧订单号(请记录到数据库中)
	Success  bool   `Field:"success"`     //兑换是否成功，状态是true和false
	Message  string `json:"errorMessage"` //消息
	Id       string `json:"bizId"`        //业务id

	Sign string `Field:"sign"` //MD5签名，详见签名规则
}

func (this *DuibaExchangeNotify) ToMap() map[string]string {
	return nil
}

//兑吧信息
type DuibaConfig struct {
	Appkey    string
	AppSecret string
}

//签名
func (this *DuibaConfig) Sign(parames map[string]string) string {
	if parames != nil {
		parames["appSecret"] = this.AppSecret
	}
	nameArray := []string{}
	for key, _ := range parames {
		nameArray = append(nameArray, key)
	}

	sort.Strings(nameArray)

	var theString = ""
	for _, name := range nameArray {
		if "sign" == name {
			continue
		}
		theString = theString + parames[name]
	}

	return codes.ToMd5Hex(theString)
}

//校验
func (this *DuibaConfig) Verify(parames map[string]string, sign string) bool {
	target := this.Sign(parames)
	if sign == target {
		return true
	}
	return false
}

func (this *DuibaConfig) OnExchangeNotify(p *DuibaExchangeNotify, f func()) string {
	//this.Verify(p.ToMap())
	return DUIBA_FAIL
}

//1、兑吧扣款接口
func (this *DuibaConfig) OnExchange(p map[string]string, f func()) *DuibaResponse {
	if this.Verify(p, p["sign"]) {
		f()
	} else {
		result := DuibaResponse{}
		result.Status = DUIBA_FAIL
		result.Message = "in developing"
		result.Credits = 0
		return &result
	}

	return nil
}
