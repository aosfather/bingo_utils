package didi

import (
	"fmt"
)

/**
滴滴评价接口
1、投诉
2、评价司机
*/
const (
	_COMMENT_DRIVER   = _DIDI_URL + "/v1/common/Comment"                 //评价司机
	_COMPLAINT_SUBMIT = _DIDI_URL + "/v1/common/Complaint/submit"        //提交投诉
	_COMPLAINT_OPTION = _DIDI_URL + "/v1/common/Complaint/getReasonList" //投诉选项
)

//评价司机
type commentRequest struct {
	didiBaseRequest
	Order   string `json:"order_id"` //订单id
	Level   string `json:"level"`    //司机评分 星级(1-5)
	Comment string `json:"comment"`  //司机评价最多40个汉字
}

func (this *commentRequest) toMap() map[string]string {
	values := this.didiBaseRequest.toMap()
	values["order_id"] = this.Order
	values["level"] = this.Level
	values["comment"] = this.Comment
	return values
}
func (this *DiDiConfig) Comment(id string, level string, comment string) error {
	r := commentRequest{}
	r.Order = id
	r.Level = level
	r.Comment = comment
	response := didiBaseResponse{}
	this.callApi(_COMMENT_DRIVER, &r.didiBaseRequest, &r, &response)
	if response.Code != ERR_SUCCESS {
		this.log.Error(response.Message)
		return fmt.Errorf("comment driver error : %d:%s", response.Code, response.Message)
	}
	return nil
}

//获取投诉选项
type complaintOption struct {
	Id      int    //投诉选项id
	Text    string //投诉选项内容
	Tag     int    `json:"p_tag_id"`   //父级分类ID
	TagName string `json:"p_tag_name"` //父级分类名称
}

type OptionResponse struct {
	didiBaseResponse
	Data []complaintOption
}

func (this *DiDiConfig) GetComplaintOption(id string) []complaintOption {
	request := orderSimpleRequest{}
	request.Order = id
	response := OptionResponse{}
	this.callGetApi(_COMPLAINT_OPTION, &request.didiBaseRequest, &request, &response)
	if response.Code == ERR_SUCCESS {
		return response.Data
	}
	this.log.Error("get option error:%d %s", response.Code, response.Message)
	return nil

}

//投诉
type complaintRequest struct {
	didiBaseRequest
	Order   string `json:"order_id"` //订单id
	Type    string `json:"type"`     //投诉选项id(从投诉选项接口返回的结果中选取)
	Content string `json:"content"`  //投诉选项外的其他投诉内容,不能多于40个汉字
}

func (this *complaintRequest) toMap() map[string]string {
	values := this.didiBaseRequest.toMap()
	values["order_id"] = this.Order
	values["type"] = this.Type
	values["content"] = this.Content
	return values
}

func (this *DiDiConfig) Complaint(id string, option string, content string) error {
	r := complaintRequest{}
	r.Order = id
	r.Type = option
	r.Content = content
	response := didiBaseResponse{}
	this.callApi(_COMPLAINT_SUBMIT, &r.didiBaseRequest, &r, &response)
	if response.Code != ERR_SUCCESS {
		this.log.Error(response.Message)
		return fmt.Errorf("complaint error:%d:%s", response.Code, response.Message)
	}
	return nil
}
