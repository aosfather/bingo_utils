package tongcheng

import (
	"fmt"
	"github.com/aosfather/bingo_utils/codes"
	utils "github.com/aosfather/bingo_utils/log"
	"sort"
)

const (
	SUCCESS              = 200
	ERROR                = 400
	SUCCESS_MESSAGE      = "SUCCESS"
	ERROR_MESSAGE        = "ERROR"
	TcSuccess            = "SUCCESS"
	TcError              = "FAILURE"
	TcEstimateUrl        = "/cardistributionapi/order/estimatePrice"     //预估价查询
	TcPriOrderCreateUrl  = "/cardistributionapi/order/create"            //快车订单
	TcTaxiOrderCreateUrl = "/cardistributionapi/taxiorder/create"        //出租车订单
	TcCancelUrl          = "/cardistributionapi/order/cancel"            //订单取消
	TcDriverLocationUrl  = "/cardistributionapi/order/getdriverlocation" //获取司机位置
	TcOrderDetailUrl     = "/cardistributionapi/order/getOrderDetail"    //订单详情
)

const (
	/*
	18	即时快车
	10	预约快车
	19	即时专车
	11	预约专车
	12	接机
	13	送机
	14	接站
	15	送站
	24	即时出租车
	65	预约出租车
	*/
	PD_TAXI             = "24"
	PD_TAXI_SCHEDULE    = "65"
	PD_SPECIAL          = 19
	PD_SPECIAL_SCHEDULE = 11
	PD_EXPRESS          = 18
	PD_EXPRESS_SCHEDULE = 10
)

type BaseResp struct {
	Status  int
	Message string
}

func (this *BaseResp) Success() {
	this.Status = SUCCESS
	this.Message = SUCCESS_MESSAGE
}
func (this *BaseResp) Error() {
	this.Status = ERROR
	this.Message = ERROR_MESSAGE
}
func (this *BaseResp) IsSuccess() bool {
	if SUCCESS == this.Status {
		return true
	}
	return false
}

type NotifyProcess func(request *TongChengNoticeRequestDto) bool
type Tongcheng struct {
	DistributorCd string
	TcAccessToken string
	Domain        string
	Notify        NotifyProcess
	logger        utils.Log
}

func (this *Tongcheng) Init() {
	if this.logger == nil {
		this.logger = utils.DefaultLog()
	}
}

type TcResp struct {
	Code   string `json:"code"`   /** 暂无 */
	Msg    string `json:"msg"`    /** 错误信息，无错误为空 */
	Status string `json:"status"` /** 状态（SUCCESS:成功，FAILURE:失败） **/
	Result string `json:"result"` /** 数据体 **/
	//	Content string `json:"content"` /** 数据体 */
	Success bool `json:"success"` /** 数据状态,需要默认为false,同程返回可能不通过一个状态记录 **/
}

func (this *TcResp) HasSuccess() {
	this.Status = TcSuccess
	this.Success = true
	return
}

// 同程返回状态定义
func (this *TcResp) isSuccess() bool {
	if TcSuccess == this.Status || true == this.Success {
		return true
	}
	return false
}

func getSortParamStr(reqMap map[string]string) string {
	//fmt.Println(len(reqMap))
	var keys []string
	var resultStr string
	for k, _ := range reqMap {
		if "sign" == k {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if "" == reqMap[k] {
			continue
		}
		if resultStr == "" {
			resultStr = k + "=" + reqMap[k]
		} else {
			resultStr = resultStr + "&" + k + "=" + reqMap[k]
		}
	}
	return resultStr
}

func makeEncryptSign(reqMap map[string]string, tcToken string) string {
	// 排序键并接拼参数
	var sortParamStr string
	var sign string
	sortParamStr = getSortParamStr(reqMap) + "&accessToken=" + tcToken

	fmt.Println("加密前sign值:" + sortParamStr)
	//加密
	sign += codes.ToMd5Hex(sortParamStr)
	fmt.Println("加密后的sign值:" + sign + "\n")

	return sign
}

func sign(reqMap map[string]string, tcToken string) string {
	return makeEncryptSign(reqMap, tcToken)
}
