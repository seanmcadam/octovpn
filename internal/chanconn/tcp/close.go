package tcp

// Close()
//func (t *TcpStruct) Close() {
//	log.Debugf("TCP Close() called")
//	close(t.Closech)
//}

func (t *TcpStruct) DoneChan() <-chan struct{} {
	return t.cx.DoneChan()
}

func (t *TcpStruct) Cancel(){
	t.cx.Cancel()
}

func (t *TcpStruct) closed() bool {
	return t.cx.Done()
}
