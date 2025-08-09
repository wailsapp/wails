package capabilities

import "encoding/json"

type Capabilities struct {
	HasNativeDrag bool
}

func (c Capabilities) AsBytes() []byte {
	// JSON encode
	result, err := json.Marshal(c)
	if err != nil {
		return []byte("{}")
	}
	return result
}
