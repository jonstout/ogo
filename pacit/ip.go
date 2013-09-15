package pacit

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

const (
	IP_ICMP     = 0x01
	IP_TCP      = 0x06
	IP_UDP      = 0x11
	IP_IPv6     = 0x29
	IP_IPv6ICMP = 0x3a
)

type IPv4 struct {
	Version        uint8 //4-bits
	IHL            uint8 //4-bits
	DSCP           uint8 //6-bits
	ECN            uint8 //2-bits
	Length         uint16
	ID             uint16
	Flags          uint16 //3-bits
	FragmentOffset uint16 //13-bits
	TTL            uint8
	Protocol       uint8
	Checksum       uint16
	NWSrc          net.IP
	NWDst          net.IP
	Options        []byte
	Data           ReadWriteMeasurer
}

func (i *IPv4) Len() (n uint16) {
	if i.Data != nil {
		return uint16(i.IHL*4) + i.Data.Len()
	}
	return uint16(i.IHL * 4)
}

func (i *IPv4) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	var verIhl uint8 = (i.Version << 4) + i.IHL
	binary.Write(buf, binary.BigEndian, verIhl)
	var dscpEcn uint8 = (i.DSCP << 2) + i.ECN
	binary.Write(buf, binary.BigEndian, dscpEcn)
	binary.Write(buf, binary.BigEndian, i.Length)
	binary.Write(buf, binary.BigEndian, i.ID)
	var flagsFrag uint16 = (i.Flags << 13) + i.FragmentOffset
	binary.Write(buf, binary.BigEndian, flagsFrag)
	binary.Write(buf, binary.BigEndian, i.TTL)
	binary.Write(buf, binary.BigEndian, i.Protocol)
	binary.Write(buf, binary.BigEndian, i.Checksum)
	binary.Write(buf, binary.BigEndian, i.NWSrc)
	binary.Write(buf, binary.BigEndian, i.NWDst)
	binary.Write(buf, binary.BigEndian, i.Options)
	if i.Data != nil {
		if n, err := buf.ReadFrom(i.Data); n == 0 {
			return int(n), err
		}
	}
	if n, err = buf.Read(b); n == 0 {
		return
	}
	return n, io.EOF
}

func (i *IPv4) ReadFrom(r io.Reader) (n int64, err error) {
	var verIhl uint8 = 0
	if err = binary.Read(r, binary.BigEndian, &verIhl); err != nil {
		return
	}
	n += 1
	i.Version = verIhl >> 4
	i.IHL = verIhl & 0x0f
	var dscpEcn uint8 = 0
	if err = binary.Read(r, binary.BigEndian, &dscpEcn); err != nil {
		return
	}
	n += 1
	i.DSCP = dscpEcn >> 2
	i.ECN = dscpEcn & 0x03
	if err = binary.Read(r, binary.BigEndian, &i.Length); err != nil {
		return
	}
	n += 2
	if err = binary.Read(r, binary.BigEndian, &i.ID); err != nil {
		return
	}
	n += 2
	var flagsFrag uint16 = 0
	if err = binary.Read(r, binary.BigEndian, &flagsFrag); err != nil {
		return
	}
	n += 2
	i.Flags = flagsFrag >> 13
	i.FragmentOffset = flagsFrag & 0x1fff
	if err = binary.Read(r, binary.BigEndian, &i.TTL); err != nil {
		return
	}
	n += 1
	if err = binary.Read(r, binary.BigEndian, &i.Protocol); err != nil {
		return
	}
	n += 1
	if err = binary.Read(r, binary.BigEndian, &i.Checksum); err != nil {
		return
	}
	n += 2
	i.NWSrc = make([]byte, 4)
	if err = binary.Read(r, binary.BigEndian, &i.NWSrc); err != nil {
		return
	}
	n += 4
	i.NWDst = make([]byte, 4)
	if err = binary.Read(r, binary.BigEndian, &i.NWDst); err != nil {
		return
	}
	n += 4
	if int(i.IHL) > 5 {
		i.Options = make([]byte, 4*(int(i.IHL)-5))
		if err = binary.Read(r, binary.BigEndian, &i.Options); err != nil {
			return
		}
		n += int64(len(i.Options))
	}
	switch i.Protocol {
	case IP_ICMP:
		trash := make([]byte, int(i.Length-20))
		binary.Read(r, binary.BigEndian, &trash)
		i.Data = new(ICMP)
		if n, err := i.Data.Read(trash); err != nil {
			return int64(n), err
		}
	case IP_UDP:
		i.Data = new(UDP)
		data := make([]byte, int(i.Length-20))
		binary.Read(r, binary.BigEndian, &data)
		if n, err := i.Data.Read(data); err != nil {
			return int64(n), err
		}
	default:
		trash := make([]byte, int(i.Length-20))
		binary.Read(r, binary.BigEndian, &trash)
	}
	n = int64(i.Length)
	return
}

func (i *IPv4) Write(b []byte) (n int, err error) {
	verIhl := b[0]
	n += 1
	i.Version = verIhl >> 4
	i.IHL = verIhl & 0x0f
	dscpEcn := b[1]
	n += 1
	i.DSCP = dscpEcn >> 2
	i.ECN = dscpEcn & 0x03
	i.Length = binary.BigEndian.Uint16(b[2:4])
	n += 2
	i.ID = binary.BigEndian.Uint16(b[4:6])
	n += 2
	flagsFrag := binary.BigEndian.Uint16(b[6:8])
	n += 2
	i.Flags = flagsFrag >> 13
	i.FragmentOffset = flagsFrag & 0x1fff
	i.TTL = b[8]
	n += 1
	i.Protocol = b[9]
	n += 1
	i.Checksum = binary.BigEndian.Uint16(b[10:12])
	n += 2
	i.NWSrc = make([]byte, 4)
	i.NWSrc = b[12:16]
	n += 4
	i.NWDst = make([]byte, 4)
	i.NWDst = b[16:20]
	n += 4
	if int(i.IHL) > 5 {
		optLen := 4 * (int(i.IHL) - 5)
		i.Options = make([]byte, optLen)
		i.Options = b[20 : 20+optLen]
		n += optLen
	}
	switch i.Protocol {
	case IP_ICMP:
		i.Data = new(ICMP)
		m, err := i.Data.Write(b[n:])
		if err != nil {
			return m, err
		}
		n += m
	case IP_UDP:
		i.Data = new(UDP)
		m, err := i.Data.Write(b[n:])
		if err != nil {
			return m, err
		}
		n += m
	default:
		panic(fmt.Sprintf("%0x\n", i.Protocol))
		//		trash := make([]byte, int(i.Length-20))
		//		binary.Read(buf, binary.BigEndian, &trash)
		n = int(i.Length)
	}
	return
}
