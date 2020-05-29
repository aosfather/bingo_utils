package reflect

import (
	"reflect"
	"strings"
)

/**
  自动压入依赖
*/

const _TAG_INJECT = "Inject" //依赖

type Object interface{}

//工厂
type Factory interface {
	CanAssignableTo(t reflect.Type) bool
	Factory(config string) interface{}
}

//服务找不到时候的处理，如果也无法处理返回nil
type OnServiceNotFoundHandler func(tag string, t reflect.Type) Object
type InjectMan struct {
	handler        OnServiceNotFoundHandler
	services       []Object
	factorys       []Factory //工厂
	serviceMapping map[string]Object
}

func (this *InjectMan) factory(tag string, t reflect.Type) Object {
	for _, f := range this.factorys {
		if f.CanAssignableTo(t) {
			return f.Factory(tag)
		}
	}
	return nil
}

func (this *InjectMan) Init(h OnServiceNotFoundHandler) {
	if h != nil {
		this.handler = h
	} else {
		this.handler = this.factory
	}

	this.serviceMapping = make(map[string]Object)
}

func (this *InjectMan) GetObjectByName(name string) Object {
	if name != "" {
		return this.serviceMapping[name]
	}

	return nil
}

func (this *InjectMan) GetObject(t reflect.Type) Object {
	for _, v := range this.services {
		if t.AssignableTo(reflect.TypeOf(v)) {
			return v
		}
	}
	return nil
}

func (this *InjectMan) AssignObject(o Object) {
	if o != nil && IsStructPtr(o) {
		t := GetRealType(o)
		o = this.GetObject(t)
	}
}

func (this *InjectMan) AddObject(o Object) {
	if o != nil && IsStructPtr(o) {
		this.addObjects(o)
	} else {
		panic("only surpport struct ptr")
	}
}

func (this *InjectMan) addObjects(o ...Object) {
	for _, obj := range o {
		this.services = append(this.services, obj)
		if f, ok := obj.(Factory); ok {
			this.factorys = append(this.factorys, f)
		}
	}

}

func (this *InjectMan) AddObjectByName(name string, o Object) {
	if o != nil && IsStructPtr(o) && name != "" {
		this.addObjects(o)
		this.serviceMapping[name] = o
	} else {
		panic("only surpport struct ptr")
	}
}

func (this *InjectMan) Inject(beforeInject func(interface{})) {
	for _, v := range this.services {
		if beforeInject != nil {
			beforeInject(v)
		}
		this.doInject(v)
	}

}

func (this *InjectMan) doInject(target Object) {
	reflectType := reflect.TypeOf(target)
	reflectValue := reflect.ValueOf(target)
	for i := 0; i < reflectValue.Elem().NumField(); i++ {
		field := reflectValue.Elem().Field(i)

		if fieldTag, ok := reflectType.Elem().Field(i).Tag.Lookup(_TAG_INJECT); ok {
			fieldType := field.Type()
			fieldName := reflectType.Elem().Field(i).Name
			//赋值，如果无法赋值则查找Set方法进行设置
			if !this.setValue(field, fieldType, fieldName, fieldTag) {
				this.setValueByMethod(reflectValue, fieldName, fieldType, fieldTag)
			}
		}

	}

	//如果存在Init方法，则调用
	initMethod := reflectValue.MethodByName("Init")
	if initMethod.IsValid() {
		if initMethod.Type().NumIn() == 0 {
			initMethod.Call([]reflect.Value{})
		}
	}
}

func (this *InjectMan) InjectObject(target Object) {
	if target != nil && IsStructPtr(target) {
		this.doInject(target)
	}
}

//调用 Setxxx方法进行设置值
func (this *InjectMan) setValueByMethod(v reflect.Value, fieldName string, ft reflect.Type, sname string) {
	methodName := "Set" + strings.ToUpper(fieldName[:1]) + fieldName[1:]
	rm := v.MethodByName(methodName)
	if rm.IsValid() {
		rm.Call([]reflect.Value{this.getReflectValue(ft, sname)})
	}

}

func (this *InjectMan) getReflectValue(t reflect.Type, sname string) reflect.Value {
	if sname != "" {
		v := this.serviceMapping[sname]
		if v == nil && this.handler != nil {
			v = this.handler(sname, t)
		}
		return reflect.ValueOf(v)
	} else {
		if t.Kind() == reflect.Map {
			return reflect.MakeMap(t)
		}

		//其它类型，轮询所有的service
		v := this.GetObject(t)
		if v == nil {
			if this.handler != nil {
				v := this.handler(sname, t)
				return reflect.ValueOf(v)
			}
		}
		return reflect.ValueOf(v)
	}
}

func (this *InjectMan) setValue(field reflect.Value, t reflect.Type, fname, sname string) bool {
	if !field.CanSet() {
		println("can not set field " + fname)
		return false
	}

	field.Set(this.getReflectValue(t, sname))
	return true

}
func isNilOrZero(v reflect.Value, t reflect.Type) bool {
	switch v.Kind() {
	default:
		return reflect.DeepEqual(v.Interface(), reflect.Zero(t).Interface())
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
}
