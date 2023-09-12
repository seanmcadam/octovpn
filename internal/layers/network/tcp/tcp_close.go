package tcp

func (t *TcpStruct) doneChan() <-chan struct{} {
	if t == nil {
		return nil
	}
	return t.cx.DoneChan()
}

func (t *TcpStruct) Cancel() {
	if t == nil {
		return
	}
	t.link.Close()
	t.cx.Cancel()
}

func (t *TcpStruct) closed() bool {
	if t == nil {
		return true
	}
	return t.cx.Done()
}
