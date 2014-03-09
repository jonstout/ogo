package arp

import (
	"encoding/hex"
	"strings"
	"testing"
)

func TestArpMarshalBinary(t *testing.T) {
	b := "   00 01 " + // HWType
		"08 00 " + // ProtoType
		"06 04 " + // HWLength ProtoLength
		"00 01 " + // Type_Request
		"00 00 00 00 00 00 " + // HWSrc
		"00 00 00 00 " + // IPSrc
		"00 00 00 00 00 00 " + // HWDst
		"00 00 00 00 "   // IPDst
	b = strings.Replace(b, " ", "", -1)

	a, _ := New(Type_Request)
	data, _ := a.MarshalBinary()
	d := hex.EncodeToString(data)
	if (len(b) != len(d)) || (b != d) {
		t.Log("Exp:", b)
		t.Log("Rec:", d)
		t.Errorf("Received length of %d, expected %d", len(d), len(b))
	}
}

func TestArpUnmarshalBinary(t *testing.T) {

}
