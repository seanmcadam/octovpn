package connmgr

import (
	"reflect"

	"gitbuh.com/seanmcadam/octovpn/interfaces"
	"github.com/seanmcadam/bufferpool"
	"github.com/seanmcadam/counter/counter16"
	"github.com/seanmcadam/counter/counterint"
	"github.com/seanmcadam/ctx"
	log "github.com/seanmcadam/loggy"
	"github.com/seanmcadam/octovpn/interfaces"
	"github.com/seanmcadam/octovpn/internal/network"
)

// connmgr
// Passed a single connection (Client or Server)
// Opens network connection
// Relays status
// Holds all open active connections
//

type Connmgr struct {
	serial           counterint.CounterStructInt
	cx               ctx.Ctx
	connectionch     chan interfaces.LayerInterface
	connections      []interfaces.LayerInterface
	connectionstatus []interfaces.LayerStatus
	recvch           chan *bufferpool.Buffer
	statusch         chan interfaces.LayerStatus
	status           interfaces.LayerStatus
}

var counter counter16

func init() {
	counter = counter16.New()
}

// Takes a configuration object
// Converts to network
// Calls new Network
// Waits on, and managed connections
// Multiplexes the connections
func New(cx ctx.Ctx, config string) (cm interfaces.LayerInterface, err error) {
	var ch chan interfaces.LayerInterface
	ch, err = network.New(config)

	if err != nil {
		return nil, err
	}

	cm = &Connmgr{
		serial:           counter.Next(),
		cx:               cx,
		connectionch:     ch,
		connections:      make([]interfaces.LayerInterface, 0),
		connectionstatus: make([]interfaces.LayerStatus, 0),
		recvch:           make(chan *bufferpool.Buffer, 5),
		statusch:         make(chan interfaces.LayerStatus, 1),
		status:           interfaces.Down,
	}

	go func(cm *Connmgr) {
		//
		// Close down the structure here
		//
		defer func() {
			cm.sendStatus(interfaces.Closed)
			close(cm.statusch)
			close(cm.recvch)
			for i := 0; i < len(cm.connections); i++ {
				cm.connections[i].Reset()
			}
		}()

		for {
			connlen := len(cm.connections)
			if connlen == 0 && cm.status != interfaces.Down {
				cm.updateStatus()
			}

			//
			// cx, connectionch, connections = 1 + 1 + n(recv) + n(status)
			//
			cases := make([]reflect.SelectCase, 2+connlen)
			cases[0] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(cm.cx.DoneChan())}
			cases[1] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(cm.connectionch)}
			for i, conn := range cm.connections {
				cases[i+2] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(conn.RecvCh())}
			}
			for i, conn := range cm.connections {
				cases[i+connlen+2] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(conn.StatusCh())}
			}

			i, recv, ok := reflect.Select(cases)

			switch i {
			case 0: // Context is closed, exit
				return
			case 1:
				if !ok {
					return
				}
				cm.connections = append(cm.connections, recv.Interface().(interfaces.LayerInterface))
				cm.connectionstatus = append(cm.connectionstatus, interfaces.Down)
				// cm.updateStatus()

			default:
				l := len(cm.connections)
				if !ok {
					if (i - 2) < l {
						cm.remove(i - 2)
					} else {
						cm.remove(i - (2 + l))
					}
				} else {
					if (i - 2) < l {
						cm.recvch <- recv.Interface().(*bufferpool.Buffer)
					} else {
						cm.connectionstatus[i-(2+l)] = recv.Interface().(interfaces.LayerStatus)
						cm.updateStatus()
					}
				}
			}
		}
	}(cm)

	return cm, nil
}

// Send()
func (cm *Connmgr) Send(b *bufferpool.Buffer) {
	for i := 0; i < len(cm.connectionstatus); i++ {
		if cm.connectionstatus[i] == interfaces.Up {
			cm.connections[i].Send(b)
			return
		}
	}

	log.Warnf("Connmgr[%s] Send Drop...", cm.serial.String())

}

// RecvCh()
func (cm *Connmgr) RecvCh() chan *bufferpool.Buffer {
	return cm.recvch
}

// StatusCh()
func (cm *Connmgr) StatusCh() chan interfaces.LayerStatus {
	return cm.statusch
}

// Reset()
func (cm *Connmgr) Reset() {
	cm.cx.Cancel()
}

// insert()
func (cm *Connmgr) insert(c interfaces.LayerInterface) {
	cm.connections = append(cm.connections, c)
}

// remove()
func (cm *Connmgr) remove(index int) {
	cm.connections = append(cm.connections[:index], cm.connections[index+1:]...)
	cm.connectionstatus = append(cm.connectionstatus[:index], cm.connectionstatus[index+1:]...)

}

// sendStatus()
func (cm *Connmgr) sendStatus(s interfaces.LayerStatus) {
	cm.statusch <- s
}

// updateStatus()
// If one connection change to up up, send up, otherwise change to down
func (cm *Connmgr) updateStatus() {
	status := interfaces.Down

	for i := 0; i < len(cm.connectionstatus); i++ {
		if cm.connectionstatus[i] == interfaces.Up {
			status = interfaces.Up
			break
		}
	}
	if cm.status != status {
		cm.status = status
		cm.sendStatus(status)
	}
}
