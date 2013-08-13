package core

import (
	"io"
	"bytes"
	"net"
	"time"
	"encoding/binary"
	//"github.com/jonstout/pacit"
)

type LinkDiscovery struct {
	SrcDPID net.HardwareAddr
	Nsec int64 /* Number of nanoseconds elapsed since Jan 1, 1970. */
}

func NewListDiscovery(srcDPID net.HardwareAddr) (d *LinkDiscovery, err error) {
	d = new(LinkDiscovery)
	d.SrcDPID = srcDPID
	d.Nsec = time.Now().UnixNano()
	return
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
	d.SrcDPID = net.HardwareAddr(b[n:n+8])
	//d.src = b[n:n+8]
	n += 8
	d.Nsec = int64(binary.BigEndian.Uint64(b[n:n+8]))
	n += 8
	return
}
