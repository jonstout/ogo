package pacit

import (
	"bytes"
	"encoding/binary"
	"io"
)

type UDP struct {
	PortSrc  uint16
	PortDst  uint16
	Length   uint16
	Checksum uint16
	Data     []byte
}

func (u *UDP) Len() (n uint16) {
	if u.Data != nil {
		return uint16(8 + len(u.Data))
	}
	return uint16(8)
}

func (u *UDP) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, u.PortSrc)
	binary.Write(buf, binary.BigEndian, u.PortDst)
	binary.Write(buf, binary.BigEndian, u.Length)
	binary.Write(buf, binary.BigEndian, u.Checksum)
	binary.Write(buf, binary.BigEndian, u.Data)
	if n, err = buf.Read(b); n == 0 {
		return
	}
	return n, io.EOF
}

func (u *UDP) ReadFrom(r io.Reader) (n int64, err error) {
	if err = binary.Read(r, binary.BigEndian, &u.PortSrc); err != nil {
		return
	}
	n += 2
	if err = binary.Read(r, binary.BigEndian, &u.PortDst); err != nil {
		return
	}
	n += 2
	if err = binary.Read(r, binary.BigEndian, &u.Length); err != nil {
		return
	}
	n += 2
	if err = binary.Read(r, binary.BigEndian, &u.Checksum); err != nil {
		return
	}
	n += 2
	if u.Length > uint16(8) {
		u.Data = make([]byte, int(u.Length-uint16(8)))
	}
	m, err := io.ReadFull(r, u.Data)
	n += int64(m)
	return
}

func (u *UDP) Write(b []byte) (n int, err error) {
	buf := bytes.NewBuffer(b)
	if err = binary.Read(buf, binary.BigEndian, &u.PortSrc); err != nil {
		return
	}
	n += 2
	if err = binary.Read(buf, binary.BigEndian, &u.PortDst); err != nil {
		return
	}
	n += 2
	if err = binary.Read(buf, binary.BigEndian, &u.Length); err != nil {
		return
	}
	n += 2
	if err = binary.Read(buf, binary.BigEndian, &u.Checksum); err != nil {
		return
	}
	n += 2
	if u.Length > 8 {
		u.Data = make([]byte, u.Length-8)
	}
	if u.Length == 0 {
		u.Data = make([]byte, buf.Len())
	}
	m, err := io.ReadFull(buf, u.Data)
	n += m
	return
}
