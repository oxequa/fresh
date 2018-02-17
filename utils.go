package fresh

import (
	"reflect"
	"strings"
)

type (
	Utils interface {
		MapDTO(interface{}, interface{})
		ToDTO(interface{}) interface{}
	}

	utils struct {
		dr      map[string]reflect.Type
		dto_tag string
	}
)

func (u *utils) MapDTO(m interface{}, d interface{}) {
	if u.dr == nil {
		u.dr = make(map[string]reflect.Type)
	}
	u.dr[reflect.TypeOf(m).String()] = reflect.TypeOf(d)
	u.dr[reflect.TypeOf(d).String()] = reflect.TypeOf(d)
}

func (u *utils) ToDTO(model interface{}) interface{} {
	fields := make(map[string]reflect.Value)
	v := reflect.ValueOf(model)
	for i := 0; i < v.NumField(); i++ {
		fieldInfo := v.Type().Field(i)
		tag := fieldInfo.Tag
		name := tag.Get(u.dto_tag)
		if name == "" {
			name = strings.ToLower(fieldInfo.Name)
		}
		fields[name] = v.Field(i)
	}
	val := reflect.New(u.dr[reflect.TypeOf(model).String()+"DTO"]).Elem()
	for i := 0; i < val.NumField(); i++ {
		fieldInfo := val.Type().Field(i)
		tag := fieldInfo.Tag
		name := tag.Get(u.dto_tag)
		for k, v := range fields {
			if k == name {
				val.Field(i).Set(v)
				break
			}
		}
	}
	return val.Interface()
}
