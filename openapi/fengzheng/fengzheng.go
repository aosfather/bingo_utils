package fengzheng

import (
	"fmt"
	"github.com/aosfather/bingo_utils/codes"
	"github.com/aosfather/bingo_utils/http"
	"sort"
)

type baseParames struct {
	App  string `json:"appkey"`
	Sign string `json:"sign"`
}

func (this *baseParames) SetBase(app, sign string) {
	this.App = app
	this.Sign = sign
}

type baseResponse struct {
	Status string `json:"status"`
	Msg    string `json:"message"`
}

type parameter interface {
	ToMap() map[string]string
	SetBase(app, sign string)
}

/**
  电话花费及流量充值接口
*/
type Fengzheng struct {
	App    string
	Secret string
	Hosts  string
}

func (this *Fengzheng) sign(parames map[string]string) string {
	nameArray := []string{}
	for key, _ := range parames {
		nameArray = append(nameArray, key)
	}

	sort.Strings(nameArray)

	var theString = this.App
	for _, name := range nameArray {

		theString = theString + parames[name]

	}
	theString += this.Secret
	return codes.ToMd5Hex(theString)

}

func (this *Fengzheng) Call(method string, p parameter, response interface{}) {
	if p != nil {
		s := this.sign(p.ToMap())
		p.SetBase(this.App, s)
		http.Post(fmt.Sprintf("%s/%s", this.Hosts, method), p, response)
	}
}
