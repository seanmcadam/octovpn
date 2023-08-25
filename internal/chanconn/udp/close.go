package udp

// Close()
//func (u *UdpStruct) Close() {
//	log.Debugf("UDP Close() called")
//	close(u.Closech)
//}

func (u *UdpStruct) DoneChan() <-chan struct{} {
	return u.cx.DoneChan()
}

func (u *UdpStruct) Cancel() {
	u.cx.Cancel()
}

func (u *UdpStruct) closed() bool {
	return u.cx.Done()
}
