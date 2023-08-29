

The Tracker captures statistics about a channel

Bandwidth (send and recieve)
Latency
Recieved Packets
Dropped Packets


SEND
Copy Send packets sent[packet count]*PacketTracker (packet, time)

RECV
Copy Recv packets recv[packet count]*PacketTracker (packet, time)


ACKNAK copy of sent packet to tracking (removed from here when ack, nak, or timeout)

ACK
copy Packet Count to ACK table

NAK
copy Packet Count to NAK table


Caculate 1 second Bytes, counts, acks, naks


Ultimatly -> 1 5 15 60 sec