// Package ofpxx defines OpenFlow message types that are version independent.
package ofpxx

import (
	"encoding/binary"
	"errors"
)

// Returns a new OpenFlow header with version field set to v1.0.
var NewOfp10Header func() Header = newHeaderGenerator(1)
// Returns a new OpenFlow header with version field set to v1.3.
var NewOfp13Header func() Header = newHeaderGenerator(4)

var messageXid uint32 = 1

func newHeaderGenerator(ver int) func() Header {
	return func() Header {
		messageXid += 1
		p := Header{uint8(ver), 0, 8, messageXid}
		return p
	}
}

type Message interface {
	Header() *Header
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

func (h *Header) MarshalBinary() (data []byte, err error) {
	data = make([]byte, 8)
	h.Version = data[0]
	h.Type = data[1]
	binary.BigEndian.PutUint16(data[2:4], h.Length)
	binary.BigEndian.PutUint32(data[4:8], h.Xid)
	return
}

func (h *Header) UnmarshalBinary(data []byte) error {
	if len(data) < 4 {
		return errors.New("The []byte is too short to unmarshel a full HelloElemHeader.")
	}
	data[0] = h.Version
	data[1] = h.Type
	h.Length = binary.BigEndian.Uint16(data[2:4])
	h.Xid = binary.BigEndian.Uint32(data[4:8])
	return nil
}

const (
	HelloElemType_VersionBitmap = iota
	)

type HelloElem interface {
	Header() *HelloElemHeader
}

type HelloElemHeader struct {
	Type uint16
	Length uint16
}

func (h *HelloElemHeader) Header() *HelloElemHeader {
	return h
}

func (h *HelloElemHeader) Len() (n uint16) {
	return 4
}

func (h *HelloElemHeader) MarshalBinary() (data []byte, err error) {
	data = make([]byte, 4)
	binary.BigEndian.PutUint16(data[:2], h.Type)
	binary.BigEndian.PutUint16(data[2:4], h.Length)
	return
}

func (h *HelloElemHeader) UnmarshalBinary(data []byte) error {
	if len(data) < 4 {
		return errors.New("The []byte is too short to unmarshal a full HelloElemHeader.")
	}
	h.Type = binary.BigEndian.Uint16(data[:2])
	h.Length = binary.BigEndian.Uint16(data[2:4])
	return nil
}

type HelloElemVersionBitmap struct {
	HelloElemHeader
	Bitmaps []uint32
}

func NewHelloElemVersionBitmap() *HelloElemVersionBitmap {
	h := new(HelloElemVersionBitmap)
	h.Type = HelloElemType_VersionBitmap
	h.Bitmaps = make([]uint32, 1)
	// 1001
	h.Bitmaps[0] = uint32(8) & uint32(1)
	h.Length = h.HelloElemHeader.Len() + uint16(len(h.Bitmaps) * 4)
	return h
}

func (h *HelloElemVersionBitmap) Len() (n uint16) {
	n = h.HelloElemHeader.Len()
	n += uint16(len(h.Bitmaps) * 4)
	return
}

func (h *HelloElemVersionBitmap) MarshalBinary() (data []byte, err error) {
	data, err = h.HelloElemHeader.MarshalBinary()
	if err != nil {
		return
	}
	
	for i, _ := range h.Bitmaps {
		bytes := make([]byte, 4)
		binary.BigEndian.PutUint32(bytes, h.Bitmaps[i])
		data = append(data, bytes...)
	}
	return
}

func (h *HelloElemVersionBitmap) UnmarshalBinary(data []byte) error {
	length := len(data)
	read := 0
	if err := h.HelloElemHeader.UnmarshalBinary(data[:4]); err != nil {
		return err
	}
	read += int(h.HelloElemHeader.Len())

	h.Bitmaps = make([]uint32, 0)
	for read < length {
		h.Bitmaps = append(h.Bitmaps, binary.BigEndian.Uint32(data[read:read+4]))
		read += 4
	}
	return nil
}

// The OFPT_HELLO message consists of an OpenFlow header plus a set of variable
// size hello elements. The version field part of the header field (see 7.1)
// must be set to the highest OpenFlow switch protocol version supported by the
// sender (see 6.3.1).  The elements field is a set of hello elements,
// containing optional data to inform the initial handshake of the connection.
// Implementations must ignore (skip) all elements of a Hello message that they
// do not support.
// The version field part of the header field (see 7.1) must be set to the 
// highest OpenFlow switch protocol version supported by the sender (see 6.3.1).
// The elements field is a set of hello elements, containing optional data to
// inform the initial handshake of the connection. Implementations must ignore
// (skip) all elements of a Hello message that they do not support.
type Hello struct {
	Header
	Elements []HelloElem
}

func NewHello(ver int) (h *Hello, err error) {
	if ver == 1 {
		h.Header = NewOfp10Header()
	} else if ver == 4 {
		h.Header = NewOfp13Header()
	} else {
		err = errors.New("New hello message with unsupported verion was attempted to be created.")
	}
	h.Elements = make([]HelloElem, 1)
	h.Elements[0] = NewHelloElemVersionBitmap()
	return
}
