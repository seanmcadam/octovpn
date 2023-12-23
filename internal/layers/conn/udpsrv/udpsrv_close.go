package udpsrv

func (u *UdpServerStruct) doneChan() <-chan struct{} {
	if u == nil {
		return nil
	}
	return u.cx.DoneChan()
}

func (u *UdpServerStruct) Cancel() {
	if u == nil {
		return
	}
	if u.udpconn != nil {
		u.udpconn.Cancel()
		u.udpconn = nil
	}
	u.link.Close()
	u.cx.Cancel()
}

func (u *UdpServerStruct) closed() bool {
	if u == nil {
		return true
	}
	return u.cx.Done()
}
