package ofp10

type ConnectionUpReactor interface {
	ConnectionUp(dpid net.HardwareAddr, features ofp10.FeaturesReply)
}

