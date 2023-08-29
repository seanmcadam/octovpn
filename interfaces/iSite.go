package interfaces


type SiteInterface interface {
	Active() bool
	Send(PacketInterface) error
	RecvChan() <-chan PacketInterface
	Reset() error
}

