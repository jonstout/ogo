package core

import (
	"bytes"
	"encoding/binary"
	"io"
	"net"
	"time"
	//"github.com/jonstout/pacit"
)

type LinkDiscovery struct {
	SrcDPID net.HardwareAddr
	Nsec    int64 /* Number of nanoseconds elapsed since Jan 1, 1970. */
	Pad []byte
	
}

func NewLinkDiscovery(srcDPID net.HardwareAddr) *LinkDiscovery {
	d := new(LinkDiscovery)
	d.SrcDPID = srcDPID
	d.Nsec = time.Now().UnixNano()
	d.Pad = make([]byte, 6)
	return d
}

func (d *LinkDiscovery) Len() uint16 {
	return 22
}

func (d *LinkDiscovery) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, d.SrcDPID)
	binary.Write(buf, binary.BigEndian, d.Nsec)
	binary.Write(buf, binary.BigEndian, d.Pad)
	n, err = buf.Read(b)
	return n, io.EOF
}

func (d *LinkDiscovery) Write(b []byte) (n int, err error) {
	buf := bytes.NewBuffer(b)
	d.SrcDPID = make([]byte, 8)
	if err = binary.Read(buf, binary.BigEndian, &d.SrcDPID); err != nil {
		return
	}
	n += 8
	if err = binary.Read(buf, binary.BigEndian, &d.Nsec); err != nil {
		return
	}
	n += 8
	d.Pad = make([]byte, 6)
	if err = binary.Read(buf, binary.BigEndian, &d.Pad); err != nil {
		return
	}
	n += 6
	return
}
