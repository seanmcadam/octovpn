package udp

func (u *UdpStruct) doneChan() <-chan struct{} {
	if u == nil {
		return nil
	}
	return u.cx.DoneChan()
}

func (u *UdpStruct) Cancel() {
	if u == nil {
		return
	}
	//u.link.Down()
	u.link.Close()
	u.cx.Cancel()
}

func (u *UdpStruct) closed() bool {
	if u == nil {
		return true
	}
	return u.cx.Done()
}
