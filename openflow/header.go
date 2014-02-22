package ofp

import "binary"

// The version specifies the OpenFlow protocol version being
// used. During the current draft phase of the OpenFlow
// Protocol, the most significant bit will be set to indicate an
// experimental version and the lower bits will indicate a
// revision number. The current version is 0x01. The final
// version for a Type 0 switch will be 0x00. The length field
// indicates the total length of the message, so no additional
// framing is used to distinguish one frame from the next.
type Header struct {
	Version uint8
	Type    uint8
	Length  uint16
	XID     uint32
}

var NewOfp10Header func() *Header = newHeaderGenerator(1)
var NewOfp13Header func() *Header = newHeaderGenerator(4)

func newHeaderGenerator(ver int) func() *Header {
	var xid uint32 = 1
	return func() *Header {
		p := new(Header)
		p.Version = uint8(1)
		p.Type = 0
		p.Length = 8
		p.XID = xid
		xid += 1

		return p
	}
}

func (h *Header) GetHeader() *Header {
	return h
}

func (h *Header) Len() (n uint16) {
	return 8
}

func (h *Header) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, h)
	n, err = buf.Read(b)
	if err != nil {
		return
	}
	return n, io.EOF
}

type HelloElemHeader struct {
	Type uint16
	Length uint16
}

func (h *HelloElemHeader) MarshelBinary(b []byte) {

}

func (h *HelloElemHeader) UnmarshelBinary(b []byte) error {
	if len(b) > 4 {
		h.Type = binary.BigEndian.Uint16(b[:2])
		h.Length = binary.BigEndian.Uint16(b[2:4])
	} else {
		return errors.New("")
	}
}

// OpenFlow ofp_hello.
// The version field part of the header field (see 7.1) must be set to the highest OpenFlow switch protocol
// version supported by the sender (see 6.3.1).
// The elements field is a set of hello elements, containing optional data to inform the initial handshake
// of the connection. Implementations must ignore (skip) all elements of a Hello message that they do not
// support.
type Hello struct {
	Header
	Elements []HelloElemHeader
}
