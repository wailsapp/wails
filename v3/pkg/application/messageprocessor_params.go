package application

import "strconv"

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
