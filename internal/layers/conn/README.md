
Conn Layer

Establish connection
Send Start
Authentication

CLI
Call connection returns suthenticated connection, or keeps trying

SRV
Wait and return each authenticated connection, parent deals with multiple connections


There are two types of connections CLI, SRV
Listen will wait for connection

Cli will send a START Packet
Srv will response with a START packet

Once receiving a START packet each side will go AUTH

Then Authentication will happen