
Sending Packets

Channel monitor
    maintains channel status (UP/DOWN,Latency,Loss,Utilization%)
    UP, Lowest Loss, Utilization, Latency

Channel Router (All, Duplicates, Single)
    Sends packets based on th monitor

Read in upper layer requests
    track packet by counter ID
    Calculate Channel Route
    track sent data
    Send Down

Receiving Packets
    Recv()
    Checks packet counter ID, if recieved already drop it
    track recv data
    Send Up


