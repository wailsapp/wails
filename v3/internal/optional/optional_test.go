package optional

import "testing"

// ----- NewVar -----

func TestNewVar_bool(t *testing.T) {
	v := NewVar(true)
	if !v.IsSet() {
		t.Error("expected IsSet true after NewVar")
	}
	if !v.Get() {
		t.Error("expected Get to return true")
	}
}

func TestNewVar_int(t *testing.T) {
	v := NewVar(42)
	if !v.IsSet() {
		t.Error("expected IsSet true")
	}
	if v.Get() != 42 {
		t.Errorf("expected 42, got %d", v.Get())
	}
}

func TestNewVar_string(t *testing.T) {
	v := NewVar("hello")
	if !v.IsSet() {
		t.Error("expected IsSet true")
	}
	if v.Get() != "hello" {
		t.Errorf("expected \"hello\", got %q", v.Get())
	}
}

func TestNewVar_zero(t *testing.T) {
	v := NewVar(0)
	if !v.IsSet() {
		t.Error("expected IsSet true even for zero value")
	}
	if v.Get() != 0 {
		t.Errorf("expected 0, got %d", v.Get())
	}
}

// ----- Var zero value (unset by default) -----

func TestVar_zeroValueUnset(t *testing.T) {
	var v Var[int]
	if v.IsSet() {
		t.Error("expected IsSet false for zero-value Var")
	}
	if v.Get() != 0 {
		t.Errorf("expected zero value 0, got %d", v.Get())
	}
}

func TestVar_zeroValueString(t *testing.T) {
	var v Var[string]
	if v.IsSet() {
		t.Error("expected IsSet false")
	}
	if v.Get() != "" {
		t.Errorf("expected empty string, got %q", v.Get())
	}
}

// ----- Set -----

func TestSet_markIsSet(t *testing.T) {
	var v Var[int]
	v.Set(7)
	if !v.IsSet() {
		t.Error("expected IsSet true after Set")
	}
	if v.Get() != 7 {
		t.Errorf("expected 7, got %d", v.Get())
	}
}

func TestSet_overwrite(t *testing.T) {
	v := NewVar(1)
	v.Set(99)
	if v.Get() != 99 {
		t.Errorf("expected 99, got %d", v.Get())
	}
	if !v.IsSet() {
		t.Error("expected IsSet true after overwrite")
	}
}

func TestSet_string(t *testing.T) {
	var v Var[string]
	v.Set("world")
	if v.Get() != "world" {
		t.Errorf("expected \"world\", got %q", v.Get())
	}
}

func TestSet_bool(t *testing.T) {
	var v Var[bool]
	v.Set(false)
	if !v.IsSet() {
		t.Error("expected IsSet true even when setting false")
	}
	if v.Get() != false {
		t.Error("expected false")
	}
}

// ----- Unset -----

func TestUnset_clearsValue(t *testing.T) {
	v := NewVar(42)
	v.Unset()
	if v.IsSet() {
		t.Error("expected IsSet false after Unset")
	}
	if v.Get() != 0 {
		t.Errorf("expected zero value 0 after Unset, got %d", v.Get())
	}
}

func TestUnset_clearsString(t *testing.T) {
	v := NewVar("data")
	v.Unset()
	if v.IsSet() {
		t.Error("expected IsSet false")
	}
	if v.Get() != "" {
		t.Errorf("expected empty string, got %q", v.Get())
	}
}

func TestUnset_thenSet(t *testing.T) {
	v := NewVar(10)
	v.Unset()
	v.Set(20)
	if !v.IsSet() {
		t.Error("expected IsSet true after re-Set")
	}
	if v.Get() != 20 {
		t.Errorf("expected 20, got %d", v.Get())
	}
}

func TestUnset_onAlreadyUnset(t *testing.T) {
	var v Var[int]
	v.Unset() // should be a no-op, no panic
	if v.IsSet() {
		t.Error("expected IsSet false")
	}
}

// ----- Get returns current value without side effects -----

func TestGet_doesNotChangeIsSet(t *testing.T) {
	v := NewVar(5)
	_ = v.Get()
	if !v.IsSet() {
		t.Error("Get should not affect IsSet")
	}
}

// ----- Bool alias -----

func TestNewBool_true(t *testing.T) {
	b := NewBool(true)
	if !b.IsSet() {
		t.Error("expected IsSet true")
	}
	if !b.Get() {
		t.Error("expected Get true")
	}
}

func TestNewBool_false(t *testing.T) {
	b := NewBool(false)
	if !b.IsSet() {
		t.Error("expected IsSet true")
	}
	if b.Get() {
		t.Error("expected Get false")
	}
}

// ----- Package-level True / False vars -----

func TestPackageVar_True(t *testing.T) {
	if !True.IsSet() {
		t.Error("True should be set")
	}
	if !True.Get() {
		t.Error("True.Get() should return true")
	}
}

func TestPackageVar_False(t *testing.T) {
	if !False.IsSet() {
		t.Error("False should be set")
	}
	if False.Get() {
		t.Error("False.Get() should return false")
	}
}

// ----- Pointer receiver methods on addressable value -----

func TestVar_pointerReceiverViaAddress(t *testing.T) {
	v := NewVar[int](100)
	ptr := &v
	ptr.Set(200)
	if ptr.Get() != 200 {
		t.Errorf("expected 200, got %d", ptr.Get())
	}
	ptr.Unset()
	if ptr.IsSet() {
		t.Error("expected unset after Unset via pointer")
	}
}
