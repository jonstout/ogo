package core

import (
	"net"
)

type LinkDisco struct {
	src net.HardwareAddr
	nsec int64 /* Number of nanoseconds elapsed since Jan 1, 1970. */
}
