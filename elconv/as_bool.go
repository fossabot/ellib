package elconv

import (
	"reflect"
	"strconv"
)

// convert value to boolean
func AsBool(v interface{}) bool {
	if v == nil {
		return false
	}
	r := reflect.ValueOf(v)
	v = AsValueRef(r).Interface()
	switch v.(type) {
	case int:
		return v.(int) > 0
	case int8:
		return v.(int8) > 0
	case int16:
		return v.(int16) > 0
	case int32:
		return v.(int32) > 0
	case int64:
		return v.(int64) > 0
	case uint:
		return v.(uint) > 0
	case uint8:
		return v.(uint8) > 0
	case uint16:
		return v.(uint16) > 0
	case uint32:
		return v.(uint32) > 0
	case uint64:
		return v.(uint64) > 0
	case float32:
		return v.(float32) > 0
	case float64:
		return v.(float64) > 0
	case []uint8:
		b, err := strconv.ParseBool(string(v.([]uint8)))
		if err == nil {
			return b
		} else {
			return len(v.([]uint8)) != 0
		}
	case string:
		b, err := strconv.ParseBool(v.(string))
		if err == nil {
			return b
		} else {
			return len(v.(string)) != 0
		}
	case bool:
		return v.(bool)
	case error:
		return false
	default:
		// check nil and empty value
		switch r.Kind() {
		case reflect.Array:
			return r.Len() != 0
		case reflect.Map, reflect.Slice:
			return !(r.IsNil() || r.Len() == 0)
		case reflect.Interface, reflect.Ptr, reflect.Chan, reflect.Func:
			return !r.IsNil()
		}
		return !reflect.DeepEqual(v, reflect.Zero(r.Type()).Interface())
	}
}
