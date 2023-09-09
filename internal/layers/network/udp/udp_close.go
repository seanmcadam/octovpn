package udp

func (u *UdpStruct) DoneChan() <-chan struct{} {
	if u == nil {
		return nil
	}
	return u.cx.DoneChan()
}

func (u *UdpStruct) Cancel() {
	if u == nil {
		return
	}
	u.cx.Cancel()
}

func (u *UdpStruct) closed() bool {
	if u == nil {
		return true
	}
	return u.cx.Done()
}
