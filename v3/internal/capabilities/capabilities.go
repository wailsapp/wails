package capabilities

import "encoding/json"

type Capabilities struct {
	HasNativeDrag bool
}

func NewCapabilities(version string) Capabilities {
	return newCapabilities(version)
}

func (c Capabilities) AsBytes() []byte {
	// JSON encode
	result, err := json.Marshal(c)
	if err != nil {
		return []byte("{}")
	}
	return result
}
