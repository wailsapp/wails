package flags

type Doctor struct {
	Common
}

func (b *Doctor) Default() *Doctor {
	return &Doctor{}
}
