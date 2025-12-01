package application

import (
	"encoding/json"
	"fmt"
)

type Args struct {
	rawData json.RawMessage
}

func (a *Args) UnmarshalJSON(data []byte) error {
	a.rawData = data

	return nil
}

func (a *Args) String() string { return string(a.rawData) }

func (a *Args) AsMap() *MapArgs {
	m := make(map[string]interface{})
	if a.rawData == nil {
		return &MapArgs{
			data: m,
		}
	}
	err := json.Unmarshal(a.rawData, &m)
	if err != nil {
		return &MapArgs{
			data: m,
		}
	}
	return &MapArgs{data: m}
}

func (a *Args) ToStruct(str any) error {
	return json.Unmarshal(a.rawData, str)
}

type MapArgs struct {
	data map[string]interface{}
}

func (a *MapArgs) String(key string) *string {
	if a == nil {
		return nil
	}
	if val := a.data[key]; val != nil {
		result := fmt.Sprintf("%v", val)
		return &result
	}
	return nil
}

func (a *MapArgs) Int(s string) *int {
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

func (a *MapArgs) UInt8(s string) *uint8 {
	if a == nil {
		return nil
	}
	if val := a.data[s]; val != nil {
		return convertNumber[uint8](val)
	}
	return nil
}

func (a *MapArgs) UInt(s string) *uint {
	if a == nil {
		return nil
	}
	if val := a.data[s]; val != nil {
		return convertNumber[uint](val)
	}
	return nil
}

func (a *MapArgs) Float64(s string) *float64 {
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

func (a *MapArgs) Bool(s string) *bool {
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
