package codes

import (
	"regexp"
)

/**
  掩码处理工具
  1、设置前后保留或替换的位数
  2、保留或者替换方式
  3、补齐，不足的时候补齐
  4、正则限定处理的范围,只支持两端方式，将字符切分成前替换的部分，和尾部不替换的部分
     例如 email：(.*)(@.*) ，只替换名称对于@xx.com不做处理，那么前后保留位数都是相对于替换部分的
*/

type Sensitive struct {
	Head    int            //头，位数
	Tail    int            //尾，位数
	Padding bool           //是否补齐
	Replace bool           //是否替换
	Pattern *regexp.Regexp //正则表达式
	Mask    rune           //掩码字符
}

//执行转换
func (this *Sensitive) Convert(source string) string {
	var s, tail, target []rune
	if this.Pattern != nil {
		groups := this.Pattern.FindSubmatch([]byte(source))
		if len(groups) > 2 {
			s = []rune(string(groups[1]))
			tail = []rune(string(groups[2]))
		}
	}

	if len(s) == 0 {
		s = []rune(source)
	}

	var c1, c2 rune
	for index, v := range s {
		if this.Replace {
			c1 = this.Mask
			c2 = v
		} else {
			c1 = v
			c2 = this.Mask
		}
		if index < this.Head || index >= (len(s)-this.Tail) {
			target = append(target, c1)
		} else {
			target = append(target, c2)
		}
	}
	//补齐处理
	if len(s) < this.Head && this.Padding {
		for i := len(s); i < this.Head; i++ {
			target = append(target, this.Mask)
		}
	}
	//处理尾部数据
	if len(tail) > 0 {
		target = append(target, tail...)
	}

	return string(target)
}
