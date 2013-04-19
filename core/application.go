package ogo

import (
	"github.com/jonstout/ogo/openflow/ofp10"
)

var Applications map[string]Application
var messageChans map[uint8][]chan ofp10.OfpMsg

type Application interface {
	InitApplication(args map[string]string)
	GetName() string
	Receive()
}

type OgoApplication struct {
	Application
	Name string
	MessageChans map[uint8]chan ofp10.OfpMsg
}

func SubscribeTo(msg uint8) chan ofp10.OfpMsg {
	ch := make(chan ofp10.OfpMsg)
	messageChans[msg] = append(messageChans[msg], ch)
	return ch
}

func (o *OgoApplication) GetName() string {
	return o.Name
}
