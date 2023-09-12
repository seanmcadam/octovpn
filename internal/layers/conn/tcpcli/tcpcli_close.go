package tcpcli

func (t *TcpClientStruct) doneChan() <-chan struct{} {
	if t == nil {
		return nil
	}
	return t.cx.DoneChan()
}

func (t *TcpClientStruct) Cancel() {
	if t == nil {
		return
	}
	if t.tcpconn != nil {
		t.tcpconn.Cancel()
		t.tcpconn = nil
	}
	t.link.Close()
	t.cx.Cancel()
}

func (t *TcpClientStruct) closed() bool {
	if t == nil {
		return true
	}
	return t.cx.Done()
}
