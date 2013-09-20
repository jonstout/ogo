package pacit

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"net"
)

const (
	ARP_REQUEST = 1
	ARP_REPLY   = 2
)

type ARP struct {
	HWType      uint16
	ProtoType   uint16
	HWLength    uint8
	ProtoLength uint8
	Operation   uint16
	HWSrc       net.HardwareAddr
	IPSrc       net.IP
	HWDst       net.HardwareAddr
	IPDst       net.IP
}

func NewARP(Operation uint16) (*ARP, error) {
	switch Operation {
	case ARP_REQUEST, ARP_REPLY:
		break
	default:
		return nil, errors.New("Invalid ARP Operation")
	}
	a := new(ARP)
	a.HWType = 1
	a.ProtoType = 0x800
	a.HWLength = 6
	a.ProtoLength = 4
	a.Operation = Operation
	a.HWSrc = net.HardwareAddr(make([]byte, 6))
	a.IPSrc = net.IP(make([]byte, 4))
	a.HWDst = net.HardwareAddr(make([]byte, 6))
	a.IPDst = net.IP(make([]byte, 4))
	return a, nil
}

func (a *ARP) Len() (n uint16) {
	n += 8
	n += uint16(a.HWLength*2 + a.ProtoLength*2)
	return
}

func (a *ARP) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, a.HWType)
	binary.Write(buf, binary.BigEndian, a.ProtoType)
	binary.Write(buf, binary.BigEndian, a.HWLength)
	binary.Write(buf, binary.BigEndian, a.ProtoLength)
	binary.Write(buf, binary.BigEndian, a.Operation)
	binary.Write(buf, binary.BigEndian, a.HWSrc)
	binary.Write(buf, binary.BigEndian, a.IPSrc)
	binary.Write(buf, binary.BigEndian, a.HWDst)
	binary.Write(buf, binary.BigEndian, a.IPDst)
	n, err = buf.Read(b)
	return n, io.EOF
}

func (a *ARP) Write(b []byte) (n int, err error) {
	a.HWType = binary.BigEndian.Uint16(b[:2])
	n += 2
	a.ProtoType = binary.BigEndian.Uint16(b[2:4])
	n += 2
	a.HWLength = b[4]
	n += 1
	a.ProtoLength = b[5]
	n += 1
	a.Operation = binary.BigEndian.Uint16(b[6:8])
	n += 2
	//a.HWSrc = make([]byte, 6)
	a.HWSrc = b[8:14]
	n += 6
	//a.IPSrc = make([]byte, 4)
	a.IPSrc = b[14:18]
	n += 4
	//a.HWDst = make([]byte, 6)
	a.HWDst = b[18:24]
	n += 6
	//a.IPDst = make([]byte, 4)
	a.IPDst = b[24:28]
	n += 4
	//n += len(b[28:])
	return
}
