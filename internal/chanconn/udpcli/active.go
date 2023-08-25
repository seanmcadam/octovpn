package udpcli

func (u *UdpClientStruct) Active() bool {
	return u.udpconn != nil
}
