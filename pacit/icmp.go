package pacit

import (
	"bytes"
	"encoding/binary"
	"io"
)

type ICMP struct {
	Type     uint8
	Code     uint8
	Checksum uint16
	Data     []byte
}

func (i *ICMP) Len() (n uint16) {
	return uint16(4 + len(i.Data))
}

func (i *ICMP) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, i.Type)
	binary.Write(buf, binary.BigEndian, i.Code)
	binary.Write(buf, binary.BigEndian, i.Checksum)
	binary.Write(buf, binary.BigEndian, i.Data)
	if n, err = buf.Read(b); n == 0 {
		return
	}
	return n, io.EOF
}

func (i *ICMP) ReadFrom(r io.Reader) (n int64, err error) {
	if err = binary.Read(r, binary.BigEndian, &i.Type); err != nil {
		return
	}
	n += 1
	if err = binary.Read(r, binary.BigEndian, &i.Code); err != nil {
		return
	}
	n += 1
	if err = binary.Read(r, binary.BigEndian, &i.Checksum); err != nil {
		return
	}
	n += 2
	return
}

func (i *ICMP) Write(b []byte) (n int, err error) {
	buf := bytes.NewBuffer(b)
	if err = binary.Read(buf, binary.BigEndian, &i.Type); err != nil {
		return
	}
	n += 1
	if err = binary.Read(buf, binary.BigEndian, &i.Code); err != nil {
		return
	}
	n += 1
	if err = binary.Read(buf, binary.BigEndian, &i.Checksum); err != nil {
		return
	}
	n += 2
	i.Data = make([]byte, len(b)-n)
	if err = binary.Read(buf, binary.BigEndian, &i.Data); err != nil {
		return
	}
	n += len(i.Data)
	return
}
