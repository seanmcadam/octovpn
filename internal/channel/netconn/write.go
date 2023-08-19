package netconn

import "fmt"

func (nc *NetConnStruct)Write(b []byte)(l int, err error){

	// Validate b size limits here

	nc.clock.Lock()
	defer nc.clock.Unlock()

	if nc.isClosed(){
		return 0, fmt.Errorf("connection closed")
	}

	p, err := NewNetworkPacket(<-nc.count,b)
	if err != nil{
		return 0, err
	}
	
	nc.sendch <- p
	return int(p.Length()), nil
}