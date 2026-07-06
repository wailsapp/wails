package types

type UUID string

func (u *UUID) Capture(values []string) error {
	//println("UUID =", values[0])
	*u = UUID(values[0])
	return nil
}
