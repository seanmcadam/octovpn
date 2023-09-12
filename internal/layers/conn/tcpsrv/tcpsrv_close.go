package tcpsrv

func (t *TcpServerStruct) doneChan() <-chan struct{} {
	if t == nil {
		return nil
	}
	return t.cx.DoneChan()
}

func (t *TcpServerStruct) Cancel() {
	if t == nil {
		return
	}
	if t.tcplistener != nil {
		t.tcplistener.Close()
		t.tcplistener = nil
	}
	if t.tcpconn != nil {
		t.tcpconn.Cancel()
		t.tcpconn = nil
	}
	t.link.Close()
	t.cx.Cancel()
}

func (t *TcpServerStruct) closed() bool {
	if t == nil {
		return true
	}
	return t.cx.Done()
}
