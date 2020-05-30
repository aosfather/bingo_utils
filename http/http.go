package http

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

func HttpsPost(theUrl string, data interface{}, result interface{}) error {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	//http cookie接口
	cookieJar, _ := cookiejar.New(nil)

	c := &http.Client{
		Jar:       cookieJar,
		Transport: tr,
	}

	content, _ := json.Marshal(data)
	fmt.Println(string(content))
	resp, err := c.Post(theUrl, "application/json;charset=utf-8", strings.NewReader(string(content)))
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(body))

	json.Unmarshal(body, result)
	return nil
}

func HttpsPostForm(theUrl string, params map[string]string, result interface{}) error {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	//http cookie接口
	cookieJar, _ := cookiejar.New(nil)

	c := &http.Client{
		Jar:       cookieJar,
		Transport: tr,
	}

	var values url.Values = make(map[string][]string)
	for key, val := range params {
		values.Set(key, val)
	}

	resp, err := c.PostForm(theUrl, values)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(body))

	json.Unmarshal(body, result)
	return nil
}

//https get
func HttpsGet(client *http.Client, url string) ([]byte, error) {

	response, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http get error : uri=%v , statusCode=%v", url, response.StatusCode)
	}
	return ioutil.ReadAll(response.Body)
}

//HTTPGet get 请求
func HTTPGet(uri string) ([]byte, error) {
	response, err := http.Get(uri)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http get error : uri=%v , statusCode=%v", uri, response.StatusCode)
	}
	return ioutil.ReadAll(response.Body)
}

func Post(theUrl string, data interface{}, result interface{}) error {
	content, _ := json.Marshal(data)
	resp, err := http.Post(theUrl, "application/json;charset=utf-8", strings.NewReader(string(content)))
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	//fmt.Println(string(body))

	json.Unmarshal(body, result)
	return nil
}

func PostByForm(theUrl string, params map[string]string, result interface{}) error {

	var values url.Values = make(map[string][]string)
	for key, val := range params {
		values.Set(key, val)
	}
	resp, err := http.PostForm(theUrl, values)
	if err != nil {
		fmt.Println(err)
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println(string(body))
	json.Unmarshal(body, result)
	return nil
}

func HTTPGetStruct(result interface{}, furl string, a ...interface{}) error {
	url := fmt.Sprintf(furl, a...)
	data, err := HTTPGet(url)
	if err == nil {
		fmt.Println(string(data))
		err = json.Unmarshal(data, result)
		fmt.Println(err)

	}
	return err
}
