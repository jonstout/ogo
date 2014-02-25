// Package ofpxx defines OpenFlow message types that are version independent.
package ofpxx

import (
	"encoding/binary"
	"errors"
)

// Returns a new OpenFlow header with version field set to v1.0.
var NewOfp10Header func() *Header = newHeaderGenerator(1)
// Returns a new OpenFlow header with version field set to v1.3.
var NewOfp13Header func() *Header = newHeaderGenerator(4)

var messageXid uint32 = 1

func newHeaderGenerator(ver int) func() *Header {
	return func() *Header {
		messageXid += 1
		p := &Header{uint8(ver), 0, 8, messageXid}
		return p
	}
}

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
	Xid     uint32
}

func (h *Header) Header() *Header {
	return h
}

func (h *Header) Len() (n uint16) {
	return 8
}

func (h *Header) MarshelBinary() (data []byte, err error) {
	data = make([]byte, 8)
	h.Version = data[0]
	h.Type = data[1]
	binary.BigEndian.PutUint16(data[2:4], h.Length)
	binary.BigEndian.PutUint32(data[4:8], h.Xid)
	return
}

func (h *Header) UnmarshelBinary(data []byte) error {
	if len(data) < 4 {
		return errors.New("The []byte is too short to unmarshel a full HelloElemHeader.")
	}
	data[0] = h.Version
	data[1] = h.Type
	h.Length = binary.BigEndian.Uint16(data[2:4])
	h.Xid = binary.BigEndian.Uint32(data[4:8])
	return nil
}

// The OFPT_HELLO message consists of an OpenFlow header plus a set of variable size hello elements.
// The version field part of the header field (see 7.1) must be set to the highest OpenFlow switch protocol
// version supported by the sender (see 6.3.1).
// The elements field is a set of hello elements, containing optional data to inform the initial handshake
// of the connection. Implementations must ignore (skip) all elements of a Hello message that they do not
// support.
type HelloElemHeader struct {
	Type uint16
	Length uint16
}

func (h *HelloElemHeader) Len() (n uint16) {
	return 4
}

func (h *HelloElemHeader) MarshelBinary() (data []byte, err error) {
	data = make([]byte, 4)
	binary.BigEndian.PutUint16(data[:2], h.Type)
	binary.BigEndian.PutUint16(data[2:4], h.Length)
	return
}

func (h *HelloElemHeader) UnmarshelBinary(data []byte) error {
	if len(data) < 4 {
		return errors.New("The []byte is too short to unmarshel a full HelloElemHeader.")
	}
	h.Type = binary.BigEndian.Uint16(data[:2])
	h.Length = binary.BigEndian.Uint16(data[2:4])
	return nil
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
