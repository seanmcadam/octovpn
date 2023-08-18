package netconn

import "fmt"

func (nc *NetConnStruct)Close()(err error){
	nc.close()
	return nil
}

func (nc *NetConnStruct)close(){

	nc.clock.Lock()
	defer nc.clock.Unlock()

	select{
	case <-nc.closech:
		return
	default:
	}

	err := nc.conn.Close()
	if err != nil {
		fmt.Printf("Error on close: %s", err)
	}

	close(nc.recvch)
	close(nc.closech)
	close(nc.sendch)

	// Drain send and recv chans here
	SEND_DRAIN:
	for{
		select {
		case s := <-nc.sendch:
			if s == nil{
				break SEND_DRAIN
			}
		}
	}
	RECV_DRAIN:
	for{
		select {
		case r := <-nc.recvch:
			if r == nil{
				break RECV_DRAIN
			}
		}
	}
}