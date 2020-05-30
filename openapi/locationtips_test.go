package openapi

import (
	"fmt"
	"testing"
)

func TestQueryLocationTips(t *testing.T) {
	fmt.Println("start")
	fmt.Println(QueryLocationTips("北京", "10001", "天安门"))
	fmt.Println(queryFromGAODE("北京", "天安门"))
	fmt.Println("baidu")
	fmt.Println(queryFromBaidu("北京", "天安门"))
}
