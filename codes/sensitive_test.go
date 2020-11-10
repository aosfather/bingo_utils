package codes

import (
	"regexp"
	"testing"
)

func TestSensitive_Convert(t *testing.T) {
	//前三后四
	s1 := Sensitive{3, 4, false, false, nil, '*'}
	t.Log(s1.Convert("1398888881050"))
	//保留前5位
	s2 := Sensitive{5, 0, false, false, nil, '*'}
	t.Log(s2.Convert("1398888881090"))
	//补齐
	s3 := Sensitive{5, 0, true, false, nil, '*'}
	t.Log(s3.Convert("12"))
	//替换模式
	s5 := Sensitive{5, 0, false, true, nil, '*'}
	t.Log(s5.Convert("1398888881090"))

	//正则
	s4 := Sensitive{5, 0, false, false, regexp.MustCompile("(.*)(@.*)"), '*'}
	t.Log(s4.Convert("qq123456789@email.com"))

	s6 := Sensitive{5, 0, true, false, regexp.MustCompile("(.*)(@.*)"), '*'}
	t.Log(s6.Convert("qq1"))

	s7 := Sensitive{5, 0, true, false, regexp.MustCompile("(.*)(@.*)"), '*'}
	t.Log(s7.Convert("qq1@email.com"))
}
