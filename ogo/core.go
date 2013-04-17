package ogo

import (
	"time"
	"github.com/ogo/openflow/ofp10"
)

type BasicApplication struct {
	OgoApplication
	echoReply chan ofp10.OfpMsg
}

func (b *BasicApplication) InitApplication(args map[string]string) {
	b.Name = "OgoCore"
	b.echoReply = SubscribeTo(ofp10.OFPT_ECHO_REPLY)
}

func (b *BasicApplication) Receive() {
	for {
		select {
		case m := <-b.echoReply:
			go b.SendEchoReply(m.DPID)
		}
	}
}

func (b *BasicApplication) SendEchoReply(dpid string) {
	req := ofp10.NewEchoRequest()
	s, ok := GetSwitch(dpid)
	if ok {
		<-time.After(time.Second * 1)
		s.Send(req)
	}
}
