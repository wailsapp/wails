package other

// OtherMethods has another method, but through a private embedded type.
type OtherMethods struct {
	otherMethodsImpl
}

type otherMethodsImpl int

// LikeThisOtherOne does nothing as well, but is different.
func (*otherMethodsImpl) LikeThisOtherOne() {}
