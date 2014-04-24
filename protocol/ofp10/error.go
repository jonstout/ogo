package ofp10

import (
	"encoding/binary"
	
	"github.com/jonstout/ogo/protocol/ofpxx"
	"github.com/jonstout/ogo/protocol/util"
)

// BEGIN: ofp10 - 5.4.4
// ofp_error_msg 1.0
type ErrorMsg struct {
	ofpxx.Header
	Code   uint16
	Data   util.Buffer
}

func NewErrorMsg() *ErrorMsg {
	e := new(ErrorMsg)
	e.Data = *util.NewBuffer(make([]byte, 0))
	return e
}

func (e *ErrorMsg) Len() (n uint16) {
	n = e.Header.Len()
	n += 2
	n += e.Data.Len()
	return
}

func (e *ErrorMsg) MarshalBinary() (data []byte, err error) {
	data = make([]byte, int(e.Len()))
	next := 0

	bytes, err := e.Header.MarshalBinary()
	copy(data[next:], bytes)
	next += len(bytes)
	binary.BigEndian.PutUint16(data[next:], e.Code)
	next += 2
	bytes, err = e.Data.MarshalBinary()
	copy(data[next:], bytes)
	next += len(bytes)
	return
}

func (e *ErrorMsg) UnmarshalBinary(data []byte) error {
	next := 0
	e.Header.UnmarshalBinary(data[next:])
	next += int(e.Header.Len())
	e.Code = binary.BigEndian.Uint16(data[next:])
	next += 2
	e.Data.UnmarshalBinary(data[next:])
	next += int(e.Data.Len())
	return nil
}

// ofp_error_type 1.0
const (
	ET_HELLO_FAILED = iota
	ET_BAD_REQUEST
	ET_BAD_ACTION
	ET_FLOW_MOD_FAILED
	ET_PORT_MOD_FAILED
	ET_QUEUE_OP_FAILED
)

// ofp_hello_failed_code 1.0
const (
	HFC_INCOMPATIBLE = iota
	HFC_EPERM
)

// ofp_bad_request_code 1.0
const (
	BRC_BAD_VERSION = iota
	BRC_BAD_TYPE
	BRC_BAD_STAT
	BRC_BAD_VENDOR

	BRC_BAD_SUBTYPE
	BRC_EPERM
	BRC_BAD_LEN
	BRC_BUFFER_EMPTY
	BRC_BUFFER_UNKNOWN
)

// ofp_bad_action_code 1.0
const (
	BAC_BAD_TYPE = iota
	BAC_BAD_LEN
	BAC_BAD_VENDOR
	BAC_BAD_VENDOR_TYPE
	BAC_BAD_OUT_PORT
	BAC_BAD_ARGUMENT
	BAC_EPERM
	BAC_TOO_MANY
	BAC_BAD_QUEUE
)

// ofp_flow_mod_failed_code 1.0
const (
	FMFC_ALL_TABLES_FULL = iota
	FMFC_OVERLAP
	FMFC_EPERM
	FMFC_BAD_EMERG_TIMEOUT
	FMFC_BAD_COMMAND
	FMFC_UNSUPPORTED
)

// ofp_port_mod_failed_code 1.0
const (
	PMFC_BAD_PORT = iota
	PMFC_BAD_HW_ADDR
)

// ofp_queue_op_failed_code 1.0
const (
	QOFC_BAD_PORT = iota
	QOFC_BAD_QUEUE
	QOFC_EPERM
)

// END: ofp10 - 5.4.4
// END: ofp10 - 5.4
