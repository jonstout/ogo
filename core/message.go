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
	src net.HardwareAddr
	nsec uint64 /* Number of nanoseconds elapsed since Jan 1, 1970. */
}

func NewListDiscovery(srcDPID net.HardwareAddr) (d *LinkDiscovery, err error) {
	d = new(LinkDiscovery)
	d.src = srcDPID
	d.nsec = uint64(time.Now().UnixNano())
	return
}

func (d *LinkDiscovery) Len() uint16 {
	return 16
}

func (d *LinkDiscovery) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, d.src)
	binary.Write(buf, binary.BigEndian, d.nsec)
	n, err = buf.Read(b)
	return n, io.EOF
}

func (d *LinkDiscovery) Write(b []byte) (n int, err error) {
	//d.src = make([]byte, 8)
	d.src = b[n:n+8]
	n += 8
	d.nsec = binary.BigEndian.Uint64(b[n:n+8])
	n += 8
	return
}
