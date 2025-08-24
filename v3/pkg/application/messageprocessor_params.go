package application

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type BodyParams []byte

// NewBodyParams creates BodyParams from an interface{} (typically from JSON unmarshaling)
func NewBodyParams(data interface{}) (BodyParams, error) {
	if data == nil {
		return nil, nil
	}
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return BodyParams(bytes), nil
}

// ToBodyParams is a convenience function to convert interface{} to BodyParams
// Usage: params, err := ToBodyParams(body.Params)
func ToBodyParams(data interface{}) (BodyParams, error) {
	return NewBodyParams(data)
}

func (qp BodyParams) getMap() (map[string]any, error) {
	if qp == nil {
		return nil, nil
	}
	var m map[string]any
	err := json.Unmarshal([]byte(qp), &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (qp BodyParams) String(key string) *string {
	m, err := qp.getMap()
	if err != nil || m == nil {
		return nil
	}
	if val, ok := m[key]; ok && val != nil {
		result := fmt.Sprintf("%v", val)
		return &result
	}
	return nil
}

func (qp BodyParams) Int(key string) *int {
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

func (qp BodyParams) UInt8(key string) *uint8 {
	val := qp.Int(key)
	if val == nil {
		return nil
	}
	intResult := *val
	if intResult < 0 {
		intResult = 0
	}
	if intResult > 255 {
		intResult = 255
	}
	result := uint8(intResult)
	return &result
}

func (qp BodyParams) UInt(key string) *uint {
	val := qp.Int(key)
	if val == nil {
		return nil
	}
	intResult := *val
	if intResult < 0 {
		intResult = 0
	}
	result := uint(intResult)
	return &result
}

func (qp BodyParams) Bool(key string) *bool {
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

func (qp BodyParams) Float64(key string) *float64 {
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

func (qp BodyParams) ToStruct(str any) error {
	if qp == nil {
		return nil
	}
	return json.Unmarshal([]byte(qp), str)
}
