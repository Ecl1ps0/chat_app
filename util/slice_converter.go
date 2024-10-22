package util

import "errors"

func ToInterfaceSlice[T any](items []T) []interface{} {
	interfaceSlice := make([]interface{}, len(items))
	for i, item := range items {
		interfaceSlice[i] = item
	}
	return interfaceSlice
}

func ToType[T any](input []interface{}) ([]T, error) {
	var result []T
	for _, v := range input {
		castedVal, ok := v.(T)
		if !ok {
			return nil, errors.New("fail to cast to type")
		}
		result = append(result, castedVal)
	}

	return result, nil
}
