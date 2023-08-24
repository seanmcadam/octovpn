package tcpcli

func (t *TcpClientStruct) Active() bool {
	return t.tcpconn != nil
}
