package json

type Marshaler interface {
	MarshalJSON() ([]byte, error)
}

type TextMarshaler interface {
	MarshalText() (text []byte, err error)
}
