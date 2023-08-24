package interfaces


type ChannelInterface interface {
	Active() bool
	Send([]byte) (error)
	Recv() ([]byte, error)
	Reset() error
	Close()
}
