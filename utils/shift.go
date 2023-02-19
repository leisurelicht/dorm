package utils

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"reflect"
)

// StringSlice2Map 将字符串slice转为map
func StringSlice2Map(s []string) (res map[string]interface{}) {
	res = make(map[string]interface{})
	for _, e := range s {
		res[e] = nil
	}
	return res
}

func DecodeByTag(input interface{}, output interface{}, tag string) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:  output,
		TagName: tag,
	})
	if err != nil {
		return err
	}
	return decoder.Decode(input)
}

// Struct2Map 将结构体转为map，tag指定的字段值
func Struct2Map(obj interface{}, tag string) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	data := make(map[string]interface{}, t.NumField())

	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).Tag.Get(tag) == "" {
			continue
		}
		data[t.Field(i).Tag.Get(tag)] = v.Field(i).Interface()
	}
	return data
}

// Struct2MapFilterByKeys 将结构体指针转为map，tag指定的字段值, keys指定要转换的字段
func Struct2MapFilterByKeys(m interface{}, tagName string, keys []string) map[string]interface{} {
	pks := StringSlice2Map(keys)
	res := make(map[string]interface{}, len(keys))

	t := reflect.TypeOf(m)
	v := reflect.ValueOf(m)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}
	for k := 0; k < t.NumField(); k++ {
		if _, ok := pks[t.Field(k).Tag.Get(tagName)]; ok {
			res[t.Field(k).Tag.Get(tagName)] = v.Field(k).Interface()
		}
	}
	return res
}

// Map2SliceE converts keys and values of map to slice in unspecified order with error.
func Map2SliceE(i interface{}) (ks interface{}, vs interface{}, err error) {
	t := reflect.TypeOf(i)
	if t.Kind() != reflect.Map {
		err = fmt.Errorf("the input %#v of type %T isn't a map", i, i)
		return
	}
	m := reflect.ValueOf(i)
	l := m.Len()
	keys := m.MapKeys()
	ksT, vsT := reflect.SliceOf(t.Key()), reflect.SliceOf(t.Elem())
	ksV, vsV := reflect.MakeSlice(ksT, 0, l), reflect.MakeSlice(vsT, 0, l)
	for _, k := range keys {
		ksV = reflect.Append(ksV, k)
		vsV = reflect.Append(vsV, m.MapIndex(k))
	}
	return ksV.Interface(), vsV.Interface(), nil
}
