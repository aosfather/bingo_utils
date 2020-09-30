package lua

import (
	"fmt"
	lua "github.com/yuin/gopher-lua"
	"reflect"
	"regexp"
	"strings"
)

// Option is a configuration that is used to create a new mapper.
type LuaOption struct {
	// Function to convert a lua table key to Go's one. This defaults to "ToUpperCamelCase".
	NameFunc func(string) string

	// Returns error if unused keys exist.
	ErrorUnused bool

	// A struct tag name for lua table keys . This defaults to "gluamapper"
	TagName string
}

var camelre = regexp.MustCompile(`_([a-z])`)

// ToUpperCamelCase is an Option.NameFunc that converts strings from snake case to upper camel case.
func ToUpperCamelCase(s string) string {
	return strings.ToUpper(string(s[0])) + camelre.ReplaceAllStringFunc(s[1:len(s)], func(s string) string { return strings.ToUpper(s[1:len(s)]) })
}

func NewLuaOption() LuaOption {
	opt := LuaOption{ToUpperCamelCase, false, "bingo"}
	return opt
}

func StringArrayToLuaTable(l *lua.LState, dic []string) *lua.LTable {
	table := l.NewTable()
	for k, v := range dic {
		l.SetTable(table, lua.LNumber(k), lua.LString(v))
	}
	return table

}

func ArrayToLuaTable(l *lua.LState, dic []interface{}) *lua.LTable {
	table := l.NewTable()
	for k, v := range dic {
		if m, ok := v.(map[string]interface{}); ok {
			l.SetTable(table, lua.LNumber(k), ToLuaTable2(l, m))
		} else {
			l.SetTable(table, lua.LNumber(k), ToLuaValue(v))
		}

	}
	return table
}
func ToLuaTable(l *lua.LState, dic map[string]string) *lua.LTable {
	table := l.NewTable()
	for k, v := range dic {
		l.SetTable(table, lua.LString(k), lua.LString(v))
	}
	return table
}

func ToLuaTable2(l *lua.LState, dic map[string]interface{}) *lua.LTable {
	table := l.NewTable()
	for k, v := range dic {
		l.SetTable(table, lua.LString(k), ToLuaValue(v))
	}
	return table
}

func GetRealTypeAndValue(obj interface{}) (interface{}, reflect.Type) {
	objT := reflect.TypeOf(obj)
	if objT.Kind() == reflect.Ptr {
		v := reflect.ValueOf(obj)
		return v.Elem().Interface(), objT.Elem()
	}

	return obj, objT
}

func ToLuaValue(v interface{}) lua.LValue {
	rv, t := GetRealTypeAndValue(v)
	switch t.Kind() {
	case reflect.Int:
		return lua.LNumber(rv.(int))
	case reflect.Int64:
		return lua.LNumber(rv.(int64))
	case reflect.Uint:
		return lua.LNumber(rv.(uint))
	case reflect.Float32:
		return lua.LNumber(rv.(float32))
	case reflect.Float64:
		return lua.LNumber(rv.(float64))
	case reflect.Bool:
		return lua.LBool(rv.(bool))
	case reflect.String:
		return lua.LString(rv.(string))
	case reflect.Slice:
		return lua.LString(rv.([]byte))
	default:
		return lua.LNil

	}
}
func ToGoStringArray(lv lua.LValue) []string {
	v := lv.(*lua.LTable)
	maxn := v.MaxN()
	if maxn == 0 {
		return []string{}
	} else {
		ret := make([]string, 0, maxn)
		for i := 1; i <= maxn; i++ {
			ret = append(ret, v.RawGetInt(i).String())
		}
		return ret
	}
}

func ToGoMap(lv lua.LValue) map[string]interface{} {
	v := lv.(*lua.LTable)
	ret := make(map[string]interface{})
	maxn := v.MaxN()
	if maxn == 0 {
		opt := NewLuaOption()
		v.ForEach(func(key, value lua.LValue) {
			ret[key.String()] = ToGoValue(value, opt)
		})
	}

	return ret

}

func ToGoValue(lv lua.LValue, opt LuaOption) interface{} {
	switch v := lv.(type) {
	case *lua.LNilType:
		return nil
	case lua.LBool:
		return bool(v)
	case lua.LString:
		return string(v)
	case lua.LNumber:
		return float64(v)
	case *lua.LTable:
		maxn := v.MaxN()
		if maxn == 0 { // table
			ret := make(map[interface{}]interface{})
			v.ForEach(func(key, value lua.LValue) {
				keystr := fmt.Sprint(ToGoValue(key, opt))
				ret[opt.NameFunc(keystr)] = ToGoValue(value, opt)
			})
			return ret
		} else { // array
			ret := make([]interface{}, 0, maxn)
			for i := 1; i <= maxn; i++ {
				ret = append(ret, ToGoValue(v.RawGetInt(i), opt))
			}
			return ret
		}
	default:
		return v
	}
}
