'''TCP

Module handles an established TCP connection, sending and receiving packets.

The recv function reads an receive buffer to determine how large the next packet is

 Error conditions
- Bad Packet
- Closed connection
- Send Error
- Recv Error

On Error the State is set to NoLink and Cancel() is called

once Cancel is called link.Closed() is called
