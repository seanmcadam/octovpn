


Chanconn resource provide the following

	Active() bool
    Checks that the link is active, and packets wll be send as opposed to dropped

	Send(*packetchan.ChanPacket) (error)
    Send a ChanPacket

	RecvChan() <-chan *packetchan.ChanPacket
    Wait for a ChanPacket

	Reset() error
    Forces the current active connection to be restarted

	Close()
    Closes the Chanconn resource and terminates the connection