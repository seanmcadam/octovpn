package interfaces


type ChannelInterface interface {
	Send([]byte) error
	Recv() ([]byte, error)
	Reset() error
	Close()
}
