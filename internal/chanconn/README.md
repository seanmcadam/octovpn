


Chanconn resource provide the following

	Active() bool
    Checks that the link is active, and packets wll be send as opposed to dropped

	Send([]byte) (error)
    Send a Byte array as a single packet

	Recv() ([]byte, error)
    Recieve a Byte array as a single packet

	Reset()
    Forces the current active connection to be restarted

	Close()
    Closes the Chanconn resource and terminates the connection