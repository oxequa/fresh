package fresh

import (
	"reflect"
	"strings"
)

type (
	Utils interface {
		RegisterDTO(interface{}, interface{})
		MapDTO(interface{}) interface{}
	}

	utils struct {
		dr      map[string]reflect.Type
		dto_tag string
	}
)

func (u *utils) RegisterDTO(m interface{}, d interface{}) {
	if u.dr == nil {
		u.dr = make(map[string]reflect.Type)
	}
	u.dr[reflect.TypeOf(m).String()] = reflect.TypeOf(m)
	u.dr[reflect.TypeOf(d).String()] = reflect.TypeOf(d)
}

func (u *utils) MapDTO(model interface{}) interface{} {
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
	modelName := reflect.TypeOf(model).String()
	if strings.HasSuffix(modelName, strings.ToUpper(u.dto_tag)) == true {
		modelName = strings.TrimSuffix(modelName, strings.ToUpper(u.dto_tag))
	} else {
		modelName = modelName + strings.ToUpper(u.dto_tag)
	}
	val := reflect.New(u.dr[modelName]).Elem()
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
