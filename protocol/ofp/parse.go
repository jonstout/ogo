package ofp

import (
	"log"
	"io"
	
	"github.com/jonstout/ogo/protocol/ofp10"
	"github.com/jonstout/ogo/protocol/ofp13"
	"github.com/jonstout/ogo/protocol/util"
)

func Parse(b []byte) (message *util.Message, err error) {
	switch b[0] {
	case 1:
		message, err = ofp10.Parse(b)
	case 4:
		message, err = ofp13.Parse(b)
	}
	return
}
