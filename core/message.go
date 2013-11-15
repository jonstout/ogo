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
}

func NewLinkDiscovery(srcDPID net.HardwareAddr) *LinkDiscovery {
	d := new(LinkDiscovery)
	d.SrcDPID = srcDPID
	d.Nsec = time.Now().UnixNano()
	return d
}

func (d *LinkDiscovery) Len() uint16 {
	return 16
}

func (d *LinkDiscovery) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, d.SrcDPID)
	binary.Write(buf, binary.BigEndian, d.Nsec)
	n, err = buf.Read(b)
	return n, io.EOF
}

func (d *LinkDiscovery) Write(b []byte) (n int, err error) {

	buf := bytes.NewBuffer(b)
	if err = binary.Read(buf, binary.BigEndian, &d.SrcDPID); err != nil {
		return
	}
	n += 8
	var tmp uint16
	if err = binary.Read(buf, binary.BigEndian, tmp); err != nil {
		return
	}
	d.Nsec = int64(tmp)
	n += 8
	return
}
