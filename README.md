# octovpn

OctoVPN is a system for managing multiple connection streams to a target VPN device. The intention is to monitor the connections for latenty as well as loss and react accordingly. The system will detect the lowest latency path and use that first, to prevent loss the system may send a copy of a packet via one or more paths.

The service can load share as well as send the same packet in multiple paths to help insure delivery

The system does not handle:
authentication
encryption
routing
packet validation
firewalling

All of these are left to the underlying OS to manage


Example:
There are 6 connections between the local and target VPN service device. The connections are:
1) Hardwired cable provider
2) Wireless connection provider
3) Satelite connection provider
4) Cellular connection (carrier 1)
5) Cellular connection (carrier 2)
6) Cellular connection (carrier 3)
Each carrier has certain performance charatcteristics, and the system can determine these charateristics, and take advantage of them.


This tool is made for mobile connectivity where solid performance is a requirement, such as video applications.  The packet path does not matter in that all of the combined system paths look like one path to the supported network.