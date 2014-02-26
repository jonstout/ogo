package util

import (
	"bytes"
)

type Message interface {
	//encoding.BinaryMarshaler
	//encoding.BinaryUnmarshaler
	MarshalBinary() (data []byte, err error)
	UnmarshalBinary(data []byte) error
	
	Len() uint16
}

type Buffer struct { bytes.Buffer }

func NewBuffer(buf []byte) *Buffer {
	b := new(Buffer)
	b.Buffer = *bytes.NewBuffer(buf)
	return b
}

func (b *Buffer) Len() uint16 {
	return uint16(b.Len())
}

func Checksum(b []byte) uint16 {
	csumcv := len(b) - 1 // checksum coverage
	s := uint32(0)
	for i := 0; i < csumcv; i += 2 {
		s += uint32(b[i+1])<<8 | uint32(b[i])
	}
	if csumcv&1 == 0 {
		s += uint32(b[csumcv])
	}
	s = s>>16 + s&0xffff
	s = s + s>>16
	s = ^s & 0xffff
	return uint16(s<<8 | s>>(16-8))
}

/*
var ErrTruncated = errors.New("incomplete packet")

type Packet struct {
	Preamble   [7]uint8
	Delimiter  uint8
	HWDst      [6]uint8
	HWSrc      [6]uint8
	VLANHeader eth.VLAN
	Ethertype  uint16
	Payload    []uint8
}*/
