package fresh

import (
	"errors"
	"reflect"
	"strings"
)

type (
	DTO interface {
		InitDTOArray(interface{}, int) []interface{}
		DTO(interface{}, interface{}) (interface{}, error)
		DTOS([]interface{}, []interface{}) ([]interface{}, error)
	}

	dto struct {
	}
)

func (d *dto) InitDTOArray(model interface{}, len int) []interface{} {
	interfaceArray := make([]interface{}, len)
	for i, _ := range interfaceArray {
		interfaceArray[i] = reflect.New(reflect.ValueOf(model).Elem().Type()).Interface()
	}
	return interfaceArray
}

func (d *dto) DTO(sourceModel interface{}, destinationModel interface{}) (interface{}, error) {
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
		name := tag.Get("dto")
		if name == "" {
			name = strings.ToLower(fieldInfo.Name)
		}
		fields[name] = sfv.Field(i)
	}
	dfv := reflect.Indirect(reflect.ValueOf(destinationModel))
	for i := 0; i < dfv.NumField(); i++ {
		fieldInfo := dfv.Type().Field(i)
		tag := fieldInfo.Tag
		name := tag.Get("dto")
		for k, v := range fields {
			if k == name {
				dfv.Field(i).Set(v)
				break
			}
		}
	}
	return destinationModel, nil
}

func (d *dto) DTOS(sourceModels []interface{}, destinationModels []interface{}) ([]interface{}, error) {
	if len(sourceModels) != len(destinationModels) {
		return nil, errors.New("source and destination len mismatch")
	}
	for i, sm := range sourceModels {
		dm, err := d.DTO(sm, destinationModels[i])
		if err != nil {
			return nil, err
		}
		destinationModels[i] = dm
	}
	return destinationModels, nil
}
