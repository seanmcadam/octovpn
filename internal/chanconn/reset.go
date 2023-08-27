package chanconn

func (cs *ChanconnStruct) Reset() error {
	return cs.conn.Reset()
}
