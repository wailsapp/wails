// This example explores exhaustively the behaviour of encoding/json
// when handling types that implement marshaler interfaces.
//
// When encoding values, encoding/json makes no distinction
// between pointer and base types:
// if the base type implements a marshaler interface,
// be it with plain receiver or pointer receiver,
// json.Marshal picks it up.
//
// json.Marshaler is always preferred to encoding.TextMarshaler,
// without any consideration for the receiver type.
//
// When encoding map keys, on the other hand,
// the map key type must implement encoding.TextMarshaler strictly,
// i.e. if the interface is implemented only for plain receivers,
// pointer keys will not be accepted.
//
// json.Marshaler is ignored in this case, i.e. it makes no difference
// whether the key type implements it or not.
//
// Decoding behaviour w.r.t. unmarshaler types mirrors encoding behaviour.
package main

import (
	"encoding/json"
	"fmt"
)

type A struct{}
type B struct{}
type C struct{}
type D struct{}
type E struct{}
type F struct{}
type G struct{}
type H struct{}

type T struct {
	A  A
	Ap *A
	B  B
	Bp *B
	C  C
	Cp *C
	D  D
	Dp *D
	E  E
	Ep *E
	F  F
	Fp *F
	G  G
	Gp *G
	H  H
	Hp *H
}

type MT struct {
	//A  map[A]bool    // error
	//Ap map[*A]bool   // error
	//B  map[B]bool    // error
	//Bp map[*B]bool   // error
	C  map[C]bool
	Cp map[*C]bool
	//D  map[D]bool    // error
	Dp map[*D]bool
	E  map[E]bool
	Ep map[*E]bool
	//F  map[F]bool    // error
	Fp map[*F]bool
	G  map[G]bool
	Gp map[*G]bool
	//H  map[H]bool    // error
	Hp map[*H]bool
}

func (A) MarshalJSON() ([]byte, error) {
	return []byte(`"This is j A"`), nil
}

func (*B) MarshalJSON() ([]byte, error) {
	return []byte(`"This is j *B"`), nil
}

func (C) MarshalText() ([]byte, error) {
	return []byte(`This is t C`), nil
}

func (*D) MarshalText() ([]byte, error) {
	return []byte(`This is t *D`), nil
}

func (E) MarshalJSON() ([]byte, error) {
	return []byte(`"This is j E"`), nil
}

func (E) MarshalText() ([]byte, error) {
	return []byte(`This is t E`), nil
}

func (F) MarshalJSON() ([]byte, error) {
	return []byte(`"This is j F"`), nil
}

func (*F) MarshalText() ([]byte, error) {
	return []byte(`This is t *F`), nil
}

func (*G) MarshalJSON() ([]byte, error) {
	return []byte(`"This is j *G"`), nil
}

func (G) MarshalText() ([]byte, error) {
	return []byte(`This is t G`), nil
}

func (*H) MarshalJSON() ([]byte, error) {
	return []byte(`"This is j *H"`), nil
}

func (*H) MarshalText() ([]byte, error) {
	return []byte(`This is t *H`), nil
}

func main() {
	t := &T{
		A{},
		&A{},
		B{},
		&B{},
		C{},
		&C{},
		D{},
		&D{},
		E{},
		&E{},
		F{},
		&F{},
		G{},
		&G{},
		H{},
		&H{},
	}
	enc, err := json.Marshal(t)
	fmt.Println(string(enc), err)

	mt := &MT{
		//map[A]bool{A{}: true},    // error
		//map[*A]bool{&A{}: true},  // error
		//map[B]bool{B{}: true},    // error
		//map[*B]bool{&B{}: true},  // error
		map[C]bool{C{}: true},
		map[*C]bool{&C{}: true},
		//map[D]bool{D{}: true},    // error
		map[*D]bool{&D{}: true},
		map[E]bool{E{}: true},
		map[*E]bool{&E{}: true},
		//map[F]bool{F{}: true},    // error
		map[*F]bool{&F{}: true},
		map[G]bool{G{}: true},
		map[*G]bool{&G{}: true},
		//map[H]bool{H{}: true},    // error
		map[*H]bool{&H{}: true},
	}
	enc, err = json.Marshal(mt)
	fmt.Println(string(enc), err)
}
