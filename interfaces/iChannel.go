package interfaces


//
// Interface between the Channel and ChanConn Layers
// All ChanConn objects are ChannleInterfaces
//
type ChannelInterface interface {
	Active() bool
	Send([]byte) (error)
	Recv() ([]byte, error)
	Reset() error
	Close()
}
