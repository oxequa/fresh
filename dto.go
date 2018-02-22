package fresh

import (
	"errors"
	"reflect"
	"strings"
)

type (
	DTO interface {
		MakeDTOArray(interface{}, int) []interface{}
		MapDTO(interface{}, interface{}) (interface{}, error)
		MapDTOArray([]interface{}, []interface{}) ([]interface{}, error)
	}

	dto struct {
		dto_tag string
	}
)

func (d *dto) MakeDTOArray(model interface{}, len int) []interface{} {
	interfaceArray := make([]interface{}, len)
	for i, _ := range interfaceArray {
		interfaceArray[i] = reflect.New(reflect.ValueOf(model).Elem().Type()).Interface()
	}
	return interfaceArray
}

func (d *dto) MapDTO(sourceModel interface{}, destinationModel interface{}) (interface{}, error) {
	if sourceModel == nil || destinationModel == nil {
		return nil, errors.New("source or destination struct is nil")
	}
	if reflect.ValueOf(destinationModel).Kind() != reflect.Ptr {
		return nil, errors.New("destination must be a pointer")
	}
	fields := make(map[string]reflect.Value)
	sfv := reflect.Indirect(reflect.ValueOf(sourceModel))
	for i := 0; i < sfv.NumField(); i++ {
		fieldInfo := sfv.Type().Field(i)
		tag := fieldInfo.Tag
		name := tag.Get(d.dto_tag)
		if name == "" {
			name = strings.ToLower(fieldInfo.Name)
		}
		fields[name] = sfv.Field(i)
	}
	dfv := reflect.Indirect(reflect.ValueOf(destinationModel))
	for i := 0; i < dfv.NumField(); i++ {
		fieldInfo := dfv.Type().Field(i)
		tag := fieldInfo.Tag
		name := tag.Get(d.dto_tag)
		for k, v := range fields {
			if k == name {
				dfv.Field(i).Set(v)
				break
			}
		}
	}
	return destinationModel, nil
}

func (d *dto) MapDTOArray(sourceModels []interface{}, destinationModels []interface{}) ([]interface{}, error) {
	if len(sourceModels) != len(destinationModels) {
		return nil, errors.New("source and destination len mismatch")
	}
	for i, sm := range sourceModels {
		dm, err := d.MapDTO(sm, destinationModels[i])
		if err != nil {
			return nil, err
		}
		destinationModels[i] = dm
	}
	return destinationModels, nil
}
