package ofp10

import (
	"encoding/hex"
	"net"
	"strings"
	"testing"
)

func TestFeaturesReplyMarshalBinary(t *testing.T) {
	b := "   01 06 00 20 00 00 00 02" + // Header
		"00 00 00 00 00 00 00 00" + // DPID
		"00 00 00 00" + // Buffers
		"00 00 00 00" + // Tables and pad
		"00 00 00 00" + // Capabilities
		"00 00 00 00"   // Actions
	b = strings.Replace(b, " ", "", -1)

	f := NewFeaturesReply()
	data, _ := f.MarshalBinary()
	d := hex.EncodeToString(data)
	if (len(b) != len(d)) || (b != d) {
		t.Log("Exp:", b)
		t.Log("Rec:", d)
		t.Errorf("Received length of %d, expected %d", len(d), len(b))
	}
}

func TestFeaturesReplyUnmarshalBinary(t *testing.T) {
	b := "   01 06 00 20 00 00 00 02" + // Header
		"01 02 03 04 05 06 07 08" + // DPID
		"00 00 00 00" + // Buffers
		"00 00 00 00" + // Tables and pad
		"00 00 00 00" + // Capabilities
		"00 00 00 00"   // Actions
	b = strings.Replace(b, " ", "", -1)
	bytes, _ := hex.DecodeString(b)

	f := NewFeaturesReply()
	f.UnmarshalBinary(bytes)
	if f.Header.Length != 32 {
		t.Logf("Got length %d, expected %d.", f.Header.Length, 32)
	}

	dpid, _ := net.ParseMAC("01:02:03:04:05:06:07:08")
	if len(f.DPID) != len(dpid) {
		t.Errorf("Got dpid length %d, expected %d.", len(f.DPID), len(dpid))
	}
	for i, _ := range f.DPID {
		if f.DPID[i] != dpid[i] {
			t.Log("Exp:", dpid)
			t.Log("Rec:", f.DPID)
			t.Error("DPID was parsed incorrectly.")
		}
	}
}
