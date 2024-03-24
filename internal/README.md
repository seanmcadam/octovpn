
Bottom Layer is the network layer
- Network
 - tcp
  - server
  - client
  - connection
 - udp
  - server
  - client
  - connection

Connmgr layer tracks the status of the lower layer
- Connmgr
 - auth
 - ping

A Site conists of set of connections to a remote site
- Site

Natting could happen here or betwen other layers
- NAT (Future)

Router manages site routing converting ip routes to site numbers
- Router

Transport is a layer between the local machine and the Router
It also handles injecting routes into the machines routing table
- Transport

The interface layer creates a Tun/Tap device on the local machine to interface with the machine or the local network
- Interface