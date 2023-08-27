


Chanconn resource provide the following

	Active() bool
    Checks that the link is active, and packets wll be send as opposed to dropped

    Send and Recieve from the layer above

	Send(*packetchan.ChanPacket) (error)
    Send a ChanPacket
    Wrap in ConnPacket and send to the connection layer

	RecvChan() <-chan *packetchan.ChanPacket
    Wait for a ChanPacket
    get a ConnPacket and Unwrap the packet
    If it is a ChanPacket send it up a layer

	Reset() error
    Forces the current active connection to be restarted

	Close()
    Closes the Chanconn resource and terminates the connection