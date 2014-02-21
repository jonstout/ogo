package ofp13

import (
	"errors"
)

func Parse(b []byte) (message *ofp.Message, err error) {
	switch buf[1] {
	default:
		err = errors.New("An unknown v1.3 packet type was received. Parse function will discard data.")
	}
	return
}
