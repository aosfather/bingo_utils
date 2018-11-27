package bingo_utils

import (
	"fmt"
	"strconv"
	"strings"
)

type BingoString string

func (this BingoString) SnakeString() string {
	data := make([]byte, 0, len(this)*2)
	j := false
	num := len(this)
	for i := 0; i < num; i++ {
		d := this[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}
func (this BingoString) Exist() bool {
	return string(this) != string(0x1E)
}
func (this BingoString) String() string {
	if this.Exist() {
		return string(this)
	}
	return ""
}

// Bool string to bool
func (this BingoString) Bool() (bool, error) {
	return strconv.ParseBool(this.String())
}

// Float32 string to float32
func (this BingoString) Float32() (float32, error) {
	v, err := strconv.ParseFloat(this.String(), 32)
	return float32(v), err
}

// Float64 string to float64
func (this BingoString) Float64() (float64, error) {
	return strconv.ParseFloat(this.String(), 64)
}

// Int string to int
func (this BingoString) Int() (int, error) {
	v, err := strconv.ParseInt(this.String(), 10, 32)
	return int(v), err
}

// Int8 string to int8
func (this BingoString) Int8() (int8, error) {
	v, err := strconv.ParseInt(this.String(), 10, 8)
	return int8(v), err
}

// Int16 string to int16
func (this BingoString) Int16() (int16, error) {
	v, err := strconv.ParseInt(this.String(), 10, 16)
	return int16(v), err
}

// Int32 string to int32
func (this BingoString) Int32() (int32, error) {
	v, err := strconv.ParseInt(this.String(), 10, 32)
	return int32(v), err
}

// Int64 string to int64
func (this BingoString) Int64() (int64, error) {
	v, err := strconv.ParseInt(this.String(), 10, 64)
	return int64(v), err
}

// Uint string to uint
func (this BingoString) Uint() (uint, error) {
	v, err := strconv.ParseUint(this.String(), 10, 32)
	return uint(v), err
}

// Uint8 string to uint8
func (this BingoString) Uint8() (uint8, error) {
	v, err := strconv.ParseUint(this.String(), 10, 8)
	return uint8(v), err
}

// Uint16 string to uint16
func (this BingoString) Uint16() (uint16, error) {
	v, err := strconv.ParseUint(this.String(), 10, 16)
	return uint16(v), err
}

// Uint32 string to uint31
func (this BingoString) Uint32() (uint32, error) {
	v, err := strconv.ParseUint(this.String(), 10, 32)
	return uint32(v), err
}

// Uint64 string to uint64
func (this BingoString) Uint64() (uint64, error) {
	v, err := strconv.ParseUint(this.String(), 10, 64)
	return uint64(v), err
}

type argInt []int

// get int by index from int slice
func (a argInt) Get(i int, args ...int) (r int) {
	if i >= 0 && i < len(a) {
		r = a[i]
	}
	if len(args) > 0 {
		r = args[0]
	}
	return
}

// ToStr interface to string
func ToStr(value interface{}, args ...int) (s string) {
	switch v := value.(type) {
	case bool:
		s = strconv.FormatBool(v)
	case float32:
		s = strconv.FormatFloat(float64(v), 'f', argInt(args).Get(0, -1), argInt(args).Get(1, 32))
	case float64:
		s = strconv.FormatFloat(v, 'f', argInt(args).Get(0, -1), argInt(args).Get(1, 64))
	case int:
		s = strconv.FormatInt(int64(v), argInt(args).Get(0, 10))
	case int8:
		s = strconv.FormatInt(int64(v), argInt(args).Get(0, 10))
	case int16:
		s = strconv.FormatInt(int64(v), argInt(args).Get(0, 10))
	case int32:
		s = strconv.FormatInt(int64(v), argInt(args).Get(0, 10))
	case int64:
		s = strconv.FormatInt(v, argInt(args).Get(0, 10))
	case uint:
		s = strconv.FormatUint(uint64(v), argInt(args).Get(0, 10))
	case uint8:
		s = strconv.FormatUint(uint64(v), argInt(args).Get(0, 10))
	case uint16:
		s = strconv.FormatUint(uint64(v), argInt(args).Get(0, 10))
	case uint32:
		s = strconv.FormatUint(uint64(v), argInt(args).Get(0, 10))
	case uint64:
		s = strconv.FormatUint(v, argInt(args).Get(0, 10))
	case string:
		s = v
	case []byte:
		s = string(v)
	default:
		s = fmt.Sprintf("%v", v)
	}
	return s
}

//String change val type to string
func String(val interface{}) string {
	if val == nil {
		return ""
	}

	switch t := val.(type) {
	case bool:
		return strconv.FormatBool(t)
	case int:
		return strconv.FormatInt(int64(t), 10)
	case int8:
		return strconv.FormatInt(int64(t), 10)
	case int16:
		return strconv.FormatInt(int64(t), 10)
	case int32:
		return strconv.FormatInt(int64(t), 10)
	case int64:
		return strconv.FormatInt(t, 10)
	case uint:
		return strconv.FormatUint(uint64(t), 10)
	case uint8:
		return strconv.FormatUint(uint64(t), 10)
	case uint16:
		return strconv.FormatUint(uint64(t), 10)
	case uint32:
		return strconv.FormatUint(uint64(t), 10)
	case uint64:
		return strconv.FormatUint(t, 10)
	case float32:
		return strconv.FormatFloat(float64(t), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64)
	case []byte:
		return string(t)
	case string:
		return t
	default:
		return fmt.Sprintf("%v", val)
	}
}

//Int64 change val type to int64
func Int64(val interface{}) int64 {
	if val == nil {
		return 0
	}

	switch t := val.(type) {
	case bool:
		if t {
			return int64(1)
		}
		return int64(0)
	case int:
		return int64(t)
	case int8:
		return int64(t)
	case int16:
		return int64(t)
	case int32:
		return int64(t)
	case int64:
		return int64(t)
	case uint:
		return int64(t)
	case uint8:
		return int64(t)
	case uint16:
		return int64(t)
	case uint32:
		return int64(t)
	case uint64:
		return int64(t)
	case float32:
		return int64(t)
	case float64:
		return int64(t)
	case []byte:
		i, _ := strconv.Atoi(string(t))
		return int64(i)
	case string:
		b, err := strconv.ParseBool(t)
		if err == nil {
			if b {
				return int64(1)
			}

			return int64(0)
		}
		i, _ := strconv.ParseFloat(t, 64)
		return int64(i)
	default:
		i, _ := strconv.ParseFloat((fmt.Sprintf("%v", t)), 64)
		return int64(i)
	}
}

//Float64 change val type to float64
func Float64(val interface{}) float64 {
	if val == nil {
		return float64(0)
	}

	switch t := val.(type) {
	case bool:
		if t {
			return float64(1)
		}

		return float64(0)
	case int:
		return float64(t)
	case int8:
		return float64(t)
	case int16:
		return float64(t)
	case int32:
		return float64(t)
	case int64:
		return float64(t)
	case uint:
		return float64(t)
	case uint8:
		return float64(t)
	case uint16:
		return float64(t)
	case uint32:
		return float64(t)
	case uint64:
		return float64(t)
	case float32:
		return float64(t)
	case float64:
		return t
	case []byte:
		i, _ := strconv.ParseFloat(string(t), 64)
		return i
	case string:
		i, _ := strconv.ParseFloat(t, 64)
		return i
	default:
		return float64(0)
	}
}

//Int change val type to int
func Int(val interface{}) int {
	return int(Int64(val))
}

//ToInt32 change val type to int32
func Int32(val interface{}) int32 {
	return int32(Int64(val))
}

//ToFloat32 change type to float32
func Float32(val interface{}) float32 {
	return float32(Float64(val))
}
