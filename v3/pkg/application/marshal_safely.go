package application

import (
	"encoding/json/v2"
	"fmt"
)

func marshalSafely(value any) (data []byte, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("json marshal panicked: %v", recovered)
			data = nil
		}
	}()

	data, err = json.Marshal(value)
	return data, err
}
