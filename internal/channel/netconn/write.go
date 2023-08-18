package netconn

import "fmt"

func (nc *NetConnStruct)Write(b []byte)(l int, err error){

	nc.clock.Lock()
	defer nc.clock.Unlock()

	if nc.isClosed(){
		return 0, fmt.Errorf("connection closed")
	}

	l = len(b)
	nc.sendch <- b
	return l, nil
}