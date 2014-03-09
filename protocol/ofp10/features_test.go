package ofp10

import (
	"encoding/hex"
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
