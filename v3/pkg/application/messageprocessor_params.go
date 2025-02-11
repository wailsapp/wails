package application

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type QueryParams map[string][]string

func (qp QueryParams) String(key string) *string {
	if qp == nil {
		return nil
	}
	values := qp[key]
	if len(values) == 0 {
		return nil
	}
	return &values[0]
}

func (qp QueryParams) Int(key string) *int {
	val := qp.String(key)
	if val == nil {
		return nil
	}
	result, err := strconv.Atoi(*val)
	if err != nil {
		return nil
	}
	return &result
}

func (qp QueryParams) UInt8(key string) *uint8 {
	val := qp.String(key)
	if val == nil {
		return nil
	}
	intResult, err := strconv.Atoi(*val)
	if err != nil {
		return nil
	}

	if intResult < 0 {
		intResult = 0
	}
	if intResult > 255 {
		intResult = 255
	}

	var result = uint8(intResult)

	return &result
}
func (qp QueryParams) UInt(key string) *uint {
	val := qp.String(key)
	if val == nil {
		return nil
	}
	intResult, err := strconv.Atoi(*val)
	if err != nil {
		return nil
	}

	if intResult < 0 {
		intResult = 0
	}
	if intResult > 255 {
		intResult = 255
	}

	var result = uint(intResult)

	return &result
}

func (qp QueryParams) Bool(key string) *bool {
	val := qp.String(key)
	if val == nil {
		return nil
	}
	result, err := strconv.ParseBool(*val)
	if err != nil {
		return nil
	}
	return &result
}

func (qp QueryParams) Float64(key string) *float64 {
	val := qp.String(key)
	if val == nil {
		return nil
	}
	result, err := strconv.ParseFloat(*val, 64)
	if err != nil {
		return nil
	}
	return &result
}

func (qp QueryParams) ToStruct(str any) error {
	args := qp["args"]
	if len(args) == 1 {
		return json.Unmarshal([]byte(args[0]), &str)
	}
	return nil
}

type Args struct {
	data map[string]any
}

func (a *Args) String(key string) *string {
	if a == nil {
		return nil
	}
	if val := a.data[key]; val != nil {
		result := fmt.Sprintf("%v", val)
		return &result
	}
	return nil
}

func (a *Args) Int(s string) *int {
	if a == nil {
		return nil
	}
	if val := a.data[s]; val != nil {
		return convertNumber[int](val)
	}
	return nil
}

func convertNumber[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64](val any) *T {
	if val == nil {
		return nil
	}
	var result T
	switch v := val.(type) {
	case T:
		result = v
	case float64:
		result = T(v)
	default:
		return nil
	}
	return &result
}

func (a *Args) UInt8(s string) *uint8 {
	if a == nil {
		return nil
	}
	if val := a.data[s]; val != nil {
		return convertNumber[uint8](val)
	}
	return nil
}

func (a *Args) UInt(s string) *uint {
	if a == nil {
		return nil
	}
	if val := a.data[s]; val != nil {
		return convertNumber[uint](val)
	}
	return nil
}

func (a *Args) Float64(s string) *float64 {
	if a == nil {
		return nil
	}
	if val := a.data[s]; val != nil {
		if result, ok := val.(float64); ok {
			return &result
		}
	}
	return nil
}

func (a *Args) Bool(s string) *bool {
	if a == nil {
		return nil
	}
	if val := a.data[s]; val != nil {
		if result, ok := val.(bool); ok {
			return &result
		}
	}
	return nil
}

func (qp QueryParams) Args() (*Args, error) {
	argData := qp["args"]
	var result = &Args{
		data: make(map[string]any),
	}
	if len(argData) == 1 {
		err := json.Unmarshal([]byte(argData[0]), &result.data)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}
