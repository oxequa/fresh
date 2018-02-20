package fresh

import (
	"errors"
	"reflect"
	"strings"
)

type (
	DTO interface {
		RegisterDTO(interface{}, interface{})
		MarshalDTO(interface{}) (interface{}, error)
		UnmarshalDTO(interface{}) (interface{}, error)

		MarshalAllDTO([]interface{}) ([]interface{}, error)
		UnmarshaAllDTO([]interface{}) ([]interface{}, error)
	}

	dto struct {
		dr      map[string]reflect.Type
		dto_tag string
	}
)

func (d *dto) RegisterDTO(source interface{}, destination interface{}) {
	if d.dr == nil {
		d.dr = make(map[string]reflect.Type)
	}
	d.dr[reflect.TypeOf(source).String()] = reflect.TypeOf(source)
	d.dr[reflect.TypeOf(destination).String()] = reflect.TypeOf(destination)
}

func (d *dto) MarshalDTO(model interface{}) (interface{}, error) {
	if model == nil {
		return nil, errors.New("no model to map")
	}
	modelName := reflect.TypeOf(model).String()
	if !strings.HasSuffix(modelName, strings.ToUpper(d.dto_tag)) {
		modelName = modelName + strings.ToUpper(d.dto_tag)
	} else {
		return nil, errors.New("the model is already a DTO")
	}
	return d.mapDTO(model, modelName)
}

func (d *dto) MarshalAllDTO(models []interface{}) (res []interface{}, err error) {
	for _, model := range models {
		dtoModel, e := d.MarshalDTO(model)
		if err != nil {
			return nil, e
		}
		res = append(res, dtoModel)
	}
	return res, nil
}

func (d *dto) UnmarshalDTO(model interface{}) (interface{}, error) {
	if model == nil {
		return nil, errors.New("no model to map")
	}
	modelName := reflect.TypeOf(model).String()
	if strings.HasSuffix(modelName, strings.ToUpper(d.dto_tag)) {
		modelName = strings.TrimSuffix(modelName, strings.ToUpper(d.dto_tag))
	} else {
		return nil, errors.New("the DTO is already mapped")
	}
	return d.mapDTO(model, modelName)
}

func (d *dto) UnmarshaAllDTO(dtoModels []interface{}) (res []interface{}, err error) {
	for _, dtoModel := range dtoModels {
		model, e := d.UnmarshalDTO(dtoModel)
		if err != nil {
			return nil, e
		}
		res = append(res, model)
	}
	return res, nil
}

func (d *dto) mapDTO(model interface{}, modelName string) (interface{}, error) {
	if _, ok := d.dr[modelName]; !ok {
		return nil, errors.New("unable to find " + d.dto_tag + " mapping")
	}
	fields := make(map[string]reflect.Value)
	v := reflect.ValueOf(model)
	for i := 0; i < v.NumField(); i++ {
		fieldInfo := v.Type().Field(i)
		tag := fieldInfo.Tag
		name := tag.Get(d.dto_tag)
		if name == "" {
			name = strings.ToLower(fieldInfo.Name)
		}
		fields[name] = v.Field(i)
	}
	val := reflect.New(d.dr[modelName]).Elem()
	for i := 0; i < val.NumField(); i++ {
		fieldInfo := val.Type().Field(i)
		tag := fieldInfo.Tag
		name := tag.Get(d.dto_tag)
		for k, v := range fields {
			if k == name {
				val.Field(i).Set(v)
				break
			}
		}
	}
	return val.Interface(), nil

}
