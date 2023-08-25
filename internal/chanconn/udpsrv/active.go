package udpsrv

func (u *UdpServerStruct) Active() bool {
	return u.udpconn != nil
}
