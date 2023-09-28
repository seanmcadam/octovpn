

Link

Updates are not passing though from one link to the next.

Need a Map of actions and reactions
Also need to map out Recv


State consists of 

Uninitalized:
    LinkStateNONE
Down: indicating no connection or activity
    LinkStateLISTEN
    LinkStateNOLINK
    LinkStateSTART
Up: indicating end to end connectivity
    LinkStateLINK
    LinkStateCHAL
    LinkStateAUTH
    LinkStateCONNECTED
Error:
    LinkStateERROR


Link maintaines a state via direct input, and remote inputs.


        Local Set State  -----------> Set
Upstream Remote Set State -> Logic -> State -> Send State downstream

Generally local and remote sets should not happen at the same time.
Local can be used to manage local state and pass downstream or to set an initial state
Remote should be used to manage one or more upstream targets through the logic


Receiving Events:

AddLink<Type>Ch(LinkStateStruct)
Adds the state channel to the local event monitoring to be run through the logic when an event is recieved


Sending Events
Individual event type lists are maintained an activated when an appropriate event is received

