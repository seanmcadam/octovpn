package interfaces

// Interface between the Channel and ChanConn Layers
// All ChanConn objects are ChannleInterfaces
type ChannelInterface interface {
	Active() bool
	Send(PacketInterface) error
	RecvChan() <-chan PacketInterface
	Reset() error
	//Stats() TrackerData
}
