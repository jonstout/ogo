package ipv4

import (
	"bytes"
	"encoding/hex"
	"net"
	"strings"
	"testing"
)

func TestIPv4MarshalBinary(t *testing.T) {
	b := "   55 " + // Version, IHL
		"00 " + // DSCP, ECN
		"00 14 " + // Length
		"00 00 " + // Id
		"00 00 " + // Flags, FragmentOffset
		"01 " + // TTL
		"06 " + // Protocol
		"00 00 " + // Checksum
		"7f 00 00 01 " + // NWSrc
		"08 08 08 08 "   // NWDst
	b = strings.Replace(b, " ", "", -1)

	ip := New()
	ip.NWSrc = net.ParseIP("127.0.0.1")
	ip.NWDst = net.ParseIP("8.8.8.8")
	ip.Version = 5
	ip.Length = 20
	ip.TTL = 1
	ip.Protocol = Type_TCP

	data, _ := ip.MarshalBinary()
	d := hex.EncodeToString(data)
	if (len(b) != len(d)) || (b != d) {
		t.Log("Exp:", b)
		t.Log("Rec:", d)
		t.Errorf("Received length of %d, expected %d", len(d), len(b))
	}
}

func TestIPv4UnmarshalBinary(t *testing.T) {
	b := "   55 " + // Version, IHL
		"00 " + // DSCP, ECN
		"00 14 " + // Length
		"00 00 " + // Id
		"00 00 " + // Flags, FragmentOffset
		"01 " + // TTL
		"06 " + // Protocol
		"09 af " + // Checksum
		"7f 00 00 01 " + // NWSrc
		"08 08 08 08 "   // NWDst
	b = strings.Replace(b, " ", "", -1)
	byte, _ := hex.DecodeString(b)

	ip := New()
	ip.UnmarshalBinary(byte)

	src := net.ParseIP("127.0.0.1").To4()
	dst := net.ParseIP("8.8.8.8").To4()
	
	if int(ip.Len()) != len(byte) {
		t.Errorf("Got length of %d, expected %d.", ip.Len(), len(byte))
	} else if ip.Version != 5 {
		t.Errorf("Got type %d, expected %d.", ip.Version, 5)
	} else if ip.Length != 20 {
		t.Errorf("Got length %d, expected %d.", ip.Length, 20)
	} else if ip.TTL != 1 {
		t.Errorf("Got ttl %d, expected %d.", ip.TTL, 1)
	} else if ip.Protocol != Type_TCP {
		t.Errorf("Got protocol %d, expected %d.", ip.Protocol, Type_TCP)
	} else if bytes.Compare(ip.NWSrc, src) != 0 {
		t.Errorf("Got nw-src %d, expected %d.", ip.NWSrc, src)
	} else if bytes.Compare(ip.NWDst, dst) != 0 {
		t.Errorf("Got nw-dst %d, expected %d.", ip.NWDst, dst)
	}
}
