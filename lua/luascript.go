package lua

import (
	"fmt"
	lua "github.com/yuin/gopher-lua"
	"log"
)

var _default_option = NewLuaOption()

type LuaHandle func(*lua.LState)
type LuaLogFunction func(string)
type LuaContext map[string]interface{}
type LuaScript struct {
	pool     *LuaPool
	function *lua.FunctionProto
	Context  LuaContext
	Log      LuaLogFunction
}

//设置lua vm pool
func (this *LuaScript) SetPool(p *LuaPool) {
	this.pool = p
}

//加载脚本文件
func (this *LuaScript) Loadfile(filename string) {
	var err error
	this.function, err = CompileByfile(filename)
	if err != nil {
		log.Println(err.Error())
	}
}

//加载脚本
func (this *LuaScript) Load(name, content string) {
	var err error
	this.function, err = CompileByString(name, content)
	if err != nil {
		log.Println(err.Error())
	}
}

func (this *LuaScript) Call(before LuaHandle, after LuaHandle) (interface{}, error) {
	var l *lua.LState
	if this.pool != nil {
		l = this.pool.Get()
		defer func() {
			this.pool.Put(l)
		}()
	}
	if l == nil {
		return nil, fmt.Errorf("no lua vm!")
	}
	lfunc := l.NewFunctionFromProto(this.function)
	l.Push(lfunc)
	this.initFunctions(l)
	if before != nil {
		before(l)
	}
	err := l.PCall(0, 1, l.NewFunction(this.errHandle))
	if after != nil {
		after(l)
	}
	return ToGoValue(l.Get(-1), NewLuaDefaultOption()), err
}

func (this *LuaScript) CallByParame(funcname string, p ...interface{}) (interface{}, error) {
	var l *lua.LState
	if this.pool != nil {
		l = this.pool.Get()
		defer func() {
			this.pool.Put(l)
		}()
	}
	if l == nil {
		return nil, fmt.Errorf("no lua vm!")
	}

	//初始化vm
	lfunc := l.NewFunctionFromProto(this.function)
	l.Push(lfunc)
	this.initFunctions(l)
	err := l.PCall(0, 1, l.NewFunction(this.errHandle))
	//获取函数名称
	fn := l.GetGlobal(funcname)
	var paramters []lua.LValue
	for _, pitem := range p {
		paramters = append(paramters, ToLuaValue(pitem))
	}

	err = l.CallByParam(lua.P{
		Fn:      fn,
		NRet:    1,
		Protect: true,
	}, paramters...)

	return ToGoValue(l.Get(-1), NewLuaDefaultOption()), err
}

//错误捕获器
func (this *LuaScript) errHandle(l *lua.LState) int {
	log.Println(l.Get(-1).String())
	return 1
}

func (this *LuaScript) luaLog(l *lua.LState) int {
	content := l.Get(-1).String()
	l.Pop(1)
	if this.Log != nil {
		this.Log(content)
	} else {
		log.Println("lualog:" + content)
	}

	return 1
}

//从context中获取参数
func (this *LuaScript) luaGetContext(l *lua.LState) int {
	key := l.Get(-1).String()
	l.Pop(1)

	var value interface{}
	if this.Context != nil {
		value = this.Context[key]
	}
	//压入结果
	l.Push(ToLuaValue(value))
	return 1
}

//从context中获取参数
func (this *LuaScript) luaSetContext(l *lua.LState) int {
	value := l.Get(-1)
	l.Pop(1)

	key := l.Get(-1).String()
	l.Pop(1)

	if this.Context != nil {
		this.Context[key] = ToGoValue(value, _default_option)
	}

	return 1
}

func (this *LuaScript) initFunctions(l *lua.LState) {
	if l == nil {
		return
	}
	l.SetGlobal("getContext", l.NewFunction(this.luaGetContext))
	l.SetGlobal("setContext", l.NewFunction(this.luaSetContext))
	l.SetGlobal("log", l.NewFunction(this.luaLog))
}
