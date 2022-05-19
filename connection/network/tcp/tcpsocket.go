package network

type TCPConn struct {
	Socket net.conn
}

func TCPOpen() {

}

func (*TCPConn) Write() (count int, e error) {
	return count, e
}
func (*TCPConn) Read() (count int, e error) {
	return count, e
}
