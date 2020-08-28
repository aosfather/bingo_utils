package lua

import (
	"fmt"
	"testing"
	"time"
)

func TestLuaScript_Load(t *testing.T) {
	script := LuaScript{}
	script.pool = NewLuaPool(10, "", nil)
	s := `

function name()
    return 'dishonesty'
end


function incCount(n)
    log("haha")
    print(n)
 	a=2+1
 	print(a)
 	person = {
	  	name = "Michel",
	  	age  = "31", -- weakly input
	  	work_place = "San Jose"
    	}
 	return person
end
return "hello"
`
	script.Load("test", s)
	logfunc := func(a string) {
		t.Log(a)
	}
	script.Log = logfunc
	for i := 0; i < 100; i++ {
		go func() {
			v, err := script.CallByParame("incCount", fmt.Sprintf("hi%d", i))
			if err != nil {
				t.Log(err.Error())
				return
			}
			t.Log(script.Call())

			t.Log(v.(map[interface{}]interface{}))
			t.Log(v)
		}()
		time.Sleep(10 * time.Microsecond)
	}
	time.Sleep(10 * time.Second)

}
