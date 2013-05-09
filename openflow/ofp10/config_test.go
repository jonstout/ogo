package ofp10

import (
	"testing"
)

func TestSetConfig(t *testing.T) {
	c := NewSetConfig()
	if c.GetHeader().Type != OFPT_SET_CONFIG {
		t.Error("Config type was:", c.GetHeader().Type, "instead of:", OFPT_SET_CONFIG)
	}
	b := make([]byte, c.Len())
	c.Read(b)

	c2 := new(OfpSwitchConfig)
	c2.Write(b)
	if *c != *c2 {
		t.Error("Encode / Decode:", c, c2)
	}
}
