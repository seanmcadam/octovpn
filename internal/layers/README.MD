
<Down> - implying no traffic is passing
    LISTEN
    NOLINK
    START
<Up> - impling traffic is passing
    LINK
    CHAL
    AUTH
    CONNECTED

Change tcpsrv from single threaded connection to accepting unlimited numers of connections

Server: (LinkNOLINK) ->
-> Start (LinkSTART)-> Send Start 
(LinkSTART)-> Send Auth
-> Auth (LinkCHAL) -> Send Response
-> Response (LinkAuth) -> SendAccept
-> Accept (LinkCONNECT)(UP)
-> Reject (LinkNOLINK)
(LinkCONNECT) open for traffic

Client: (LinkNOLINK)
Send Start 
-> Start (LinkSTART)
(LinkSTART) -> Send Auth
-> Auth (LinkCHAL) -> Send Response
-> Response (LinkAuth) -> SendAccept
-> Accept (LinkCONNECT)(UP)
-> Reject (LinkNOLINK)
(LinkCONNECT) open for traffic
