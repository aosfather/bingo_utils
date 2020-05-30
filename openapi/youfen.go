package openapi

/**
优分
1、实现账户加载
2、实现接口访问
3、实现结果返回及转换
*/
import (
	"crypto/tls"
	"fmt"
	utils "github.com/aosfather/bingo_utils/http"
	"net/http"
)

const (
	_the_service_url = "https://api.acedata.com.cn:2443/%s?account=%s&%s"
)

type Youfen struct {
	account string
	client  *http.Client
}

func (this *Youfen) Init(account string) {
	this.account = account
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	this.client = &http.Client{Transport: tr}
}

func (this *Youfen) Call(url string, paramters string) string {
	theurl := fmt.Sprintf(_the_service_url, url, this.account, paramters)
	reponse, err := utils.HttpsGet(this.client, theurl)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	return string(reponse)
}
