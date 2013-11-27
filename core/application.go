package core

import (
	"github.com/jonstout/ogo/openflow/ofp10"
)

var Applications []InstanceGen
var messageChans map[uint8][]chan ofp10.Msg

type Application interface {
	NewInstance() interface{}
}

func SubscribeTo(msg uint8) chan ofp10.Msg {
	ch := make(chan ofp10.Msg)
	messageChans[msg] = append(messageChans[msg], ch)
	return ch
}
