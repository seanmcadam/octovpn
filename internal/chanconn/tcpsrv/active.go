package tcpsrv

func (t *TcpServerStruct) Active() bool {
	return t.tcpconn != nil
}
