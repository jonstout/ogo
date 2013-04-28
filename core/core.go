package ogo

import (
	"time"
	"github.com/jonstout/ogo/openflow/ofp10"
)

type BasicApplication struct {
	OgoApplication
	echoRequest chan ofp10.OfpMsg
}

func (b *BasicApplication) InitApplication(args map[string]string) {
	b.Name = "OgoCore"
	b.echoRequest = SubscribeTo(ofp10.OFPT_ECHO_REQUEST)
}

func (b *BasicApplication) Receive() {
	for {
		select {
		case m := <-b.echoRequest:
			go b.SendEchoReply(m.DPID)
		}
	}
}

func (b *BasicApplication) SendEchoReply(dpid string) {
	req := ofp10.NewEchoReply()
	s, ok := GetSwitch(dpid)
	if ok {
		<-time.After(time.Second * 1)
		s.Send(req)
	}
}
