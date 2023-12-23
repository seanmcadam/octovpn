package udpcli

func (u *UdpClientStruct) doneChan() <-chan struct{} {
	if u == nil {
		return nil
	}
	return u.cx.DoneChan()
}

func (u *UdpClientStruct) Cancel() {
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

func (u *UdpClientStruct) closed() bool {
	if u == nil {
		return true
	}
	return u.cx.Done()
}
