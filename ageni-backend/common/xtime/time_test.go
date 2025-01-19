package xtime

import (
	"testing"
)

func TestDuration_UnmarshalText(t *testing.T) {
	var a Duration
	err := a.UnmarshalText([]byte("4s"))
	t.Log(a)
	t.Log(err)
}

func TestTime_Scan(t *testing.T) {
	t.Log(nil)
}

func TestTime_Value(t *testing.T) {
	t.Log(nil)
}
