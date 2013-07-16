package core

/*
In general. Each switch should keep track of any switches
connected to itself.
*/
type PhyLink struct {
	port int
	latency int
	bandwidth int
}
