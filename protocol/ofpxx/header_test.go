package ofpxx

import (
	"encoding/hex"
	"strings"
	"testing"
)

func TestHelloMarshalBinary(t *testing.T) {
	b := "   01 00 00 08 00 00 00 03 " + // Header
		"00 01 00 08 " + // Element Header
		"00 00 00 09 " // Bitmap = 1001
	b = strings.Replace(b, " ", "", -1)

	h, _ := NewHello(1)
	data, _ := h.MarshalBinary()
	d := hex.EncodeToString(data)
	if (len(b) != len(d)) || (b != d) {
		t.Log("Exp:", b)
		t.Log("Rec:", d)
		t.Errorf("Received length of %d, expected %d", len(d), len(b))
	}
}

func TestHelloUnmarshalBinary(t *testing.T) {
	s := "   01 00 00 08 00 00 00 03 " + // Header
		"00 01 00 08 " + // Element Header
		"00 00 00 09 " // Bitmap = 1001
	s = strings.Replace(s, " ", "", -1)
	bytes, _ := hex.DecodeString(s)

	h, _ := NewHello(1)
	h.UnmarshalBinary(bytes)
	
	if int(h.Len()) != len(bytes) {
		t.Errorf("Got length of %d, expected %d.", h.Len(), len(bytes))
	} else if h.Version != 1 {
		t.Errorf("Got version %d, expected %d.", h.Version, 1)
	} else if h.Type != 0 {
		t.Errorf("Got type %d, expected %d.", h.Type, 0)
	} else if len(h.Elements) != 1 {
		t.Errorf("Got %d elements, expected %d elements.", len(h.Elements), 1)
	}
	
	v, ok := h.Elements[0].(*HelloElemVersionBitmap)
	if !ok {
		t.Errorf("Got wrong HelloElem type.")
	} else if len(v.Bitmaps) != 1 {
		t.Errorf("Got %d elements, expected %d elements.", len(v.Bitmaps), 1)
	} else if v.Bitmaps[0] != (uint32(8) | uint32(1)) {
		t.Errorf("Got %d bitmap, expected %d.", v.Bitmaps[0], (uint32(8) | uint32(1)))
	}
}
