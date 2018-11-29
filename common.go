package bingo_utils

import "strings"

type Object interface{}

//选项，格式 name,opt,opt1,opt2
type TagOptions []string

// Has returns true if the given optiton is available in TagOptions
func (t TagOptions) Has(opt string) bool {
	for index, tagOpt := range t {
		if index == 0 {
			continue //跳过name
		}
		if tagOpt == opt {
			return true
		}
	}

	return false
}

//将选项转换成name，opt，opt1，opt2格式
func (t *TagOptions) ToString() string {
	return strings.Join([]string(*t), ",")
}

// parseTag splits a struct field's tag into its name and a list of options
// which comes after a name. A tag is in the form of: "name,option1,option2".
// The name can be neglectected.
func (t TagOptions) ParseTag(tag string) TagOptions {
	// tag is one of followings:
	// ""
	// "name"
	// "name,opt"
	// "name,opt,opt2"
	// ",opt"
	res := strings.Split(tag, ",")
	return res
}
