package ofp

import (
	"log"
	"io"
	
	"github.com/jonstout/ogo/ofp/ofp10"
	"github.com/jonstout/ogo/ofp/ofp13"
)

// OpenFlow message header.
type Header struct {

}

type Message interface {
	io.ReadWriter
	
	GetHeader() Header
	Len() uint16
}

func Parse(b []byte) (message *Message, err error) {
	switch b[0] {
	case 1:
		message, err = ofp10.Parse(b)
	case 4:
		message, err = ofp13.Parse(b)
	}
	
	// Log all message parsing errors.
	if err != nil {
		log.Print(err)
		return
	}
}
