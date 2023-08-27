package chanconn

func (cs *ChanconnStruct) Active() bool {
	return cs.conn.Active()
}
