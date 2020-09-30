package lua

import (
	"bufio"
	"fmt"
	lua "github.com/yuin/gopher-lua"
	"github.com/yuin/gopher-lua/parse"
	"io"
	"os"
	"strings"
	"sync"
)

func SetLuaPath(p ...string) {
	luapath := "/?.lua"
	for _, path := range p {
		luapath = luapath + ";" + path + "/?.lua"
	}
	lua.LuaPathDefault = luapath
}

func AddLuaPath(p ...string) {
	luapath := lua.LuaPathDefault
	for _, path := range p {
		luapath = luapath + ";" + path + "/?.lua"
	}
	lua.LuaPathDefault = luapath
}

func CompileByfile(filename string) (*lua.FunctionProto, error) {
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(file)
	return compile(filename, reader)
}

func CompileByString(name string, content string) (*lua.FunctionProto, error) {
	return compile(name, strings.NewReader(content))
}

func compile(name string, reader io.Reader) (*lua.FunctionProto, error) {
	chunk, err := parse.Parse(reader, name)
	if err != nil {
		return nil, err
	}
	proto, err := lua.Compile(chunk, name)
	if err != nil {
		return nil, err
	}
	return proto, nil
}

//构建池
func NewLuaPool(max int, name string, lib map[string]lua.LGFunction) *LuaPool {
	pool := LuaPool{}
	pool.Init(max, name, lib)
	return &pool
}

type LuaLib struct {
	ExportName      string
	ExportFunctions map[string]lua.LGFunction
}

func (this *LuaLib) load(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), this.ExportFunctions)
	//设为只读，防止被串改
	L.Push(SetReadOnly(L, mod))
	return 1
}

//lua虚拟机池
type LuaPool struct {
	max        int
	size       int
	m          sync.Mutex
	saved      []*lua.LState
	exportName string
	exports    map[string]lua.LGFunction
	libs       []*LuaLib
}

func (this *LuaPool) InitByLibs(max int, name string, lib ...*LuaLib) {
	this.max = max
	this.saved = make([]*lua.LState, 0, max/2)
	this.libs = lib
}

func (this *LuaPool) Init(max int, name string, lib map[string]lua.LGFunction) {
	this.max = max
	this.saved = make([]*lua.LState, 0, max/2)
	this.exportName = name
	this.exports = lib
}

func (this *LuaPool) load(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), this.exports)
	//设为只读，防止被串改
	L.Push(SetReadOnly(L, mod))
	return 1
}

func (this *LuaPool) Get() *lua.LState {
	this.m.Lock()
	defer this.m.Unlock()
	if this.size < this.max {
		n := len(this.saved)
		if n == 0 {
			this.size++
			return this.new()
		}
		x := this.saved[n-1]
		this.saved = this.saved[0 : n-1]
		this.size++
		return x
	}
	return nil

}

func (this *LuaPool) new() *lua.LState {
	L := lua.NewState()
	if this.libs == nil || len(this.libs) == 0 {
		L.PreloadModule(this.exportName, this.load)
	} else {
		for _, lib := range this.libs {
			L.PreloadModule(lib.ExportName, lib.load)
		}
	}

	return L
}

//归还
func (this *LuaPool) Put(L *lua.LState) int {
	if L == nil {
		return this.size
	}
	this.m.Lock()
	defer this.m.Unlock()
	this.saved = append(this.saved, L)
	this.size--
	return this.size
}

func (this *LuaPool) Shutdown() {
	for _, L := range this.saved {
		L.Close()
	}
}

// 设置表为只读
func SetReadOnly(l *lua.LState, table *lua.LTable) *lua.LUserData {
	ud := l.NewUserData()
	mt := l.NewTable()
	// 设置表中域的指向为 table
	l.SetField(mt, "__index", table)
	// 限制对表的更新操作
	l.SetField(mt, "__newindex", l.NewFunction(func(state *lua.LState) int {
		state.RaiseError("not allow to modify table")
		return 0
	}))
	ud.Metatable = mt
	return ud
}

//检查是否使用了全局变量
//涉及到全局变量的指令有两条：GETGLOBAL（Opcode 5）和 SETGLOBAL（Opcode 7）
func CheckGlobal(proto *lua.FunctionProto) error {
	for _, code := range proto.Code {
		switch opGetOpCode(code) {
		case lua.OP_GETGLOBAL:
			return fmt.Errorf("not allow to access global")
		case lua.OP_SETGLOBAL:
			return fmt.Errorf("not allow to set global")
		}
	}
	// 对嵌套函数进行全局变量的检查
	for _, nestedProto := range proto.FunctionPrototypes {
		if err := CheckGlobal(nestedProto); err != nil {
			return err
		}
	}
	return nil
}

// 获取对应指令的 OpCode
func opGetOpCode(inst uint32) int {
	return int(inst >> 26)
}
