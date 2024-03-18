package connmgr

import (
	"net"
	"reflect"

	"github.com/seanmcadam/bufferpool"
	"github.com/seanmcadam/counter"
	"github.com/seanmcadam/counter/counterint"
	"github.com/seanmcadam/ctx"
	log "github.com/seanmcadam/loggy"
	"github.com/seanmcadam/octovpn/common"
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
	cx               *ctx.Ctx
	connectionch     chan interfaces.LayerInterface
	connections      []interfaces.LayerInterface
	connectionstatus []common.LayerStatus
	recvch           chan *bufferpool.Buffer
	statusch         chan common.LayerStatus
	status           common.LayerStatus
}

var cnt counterint.CounterStructInt

func init() {
	cnt = counter.New(ctx.New(), counter.BIT16)
}

// Takes a configuration object
// Converts to network
// Calls new Network
// Waits on, and managed connections
// Multiplexes the connections
func New(cx *ctx.Ctx, server bool, addr net.Addr) (cm interfaces.LayerInterface, err error) {
	var ch chan interfaces.LayerInterface

	if server {
		ch, err = network.Server(cx, addr)
	} else {
		ch, err = network.Client(cx, addr)
	}

	if err != nil {
		return nil, err
	}

	cm = &Connmgr{
		serial:           cnt,
		cx:               cx,
		connectionch:     ch,
		connections:      make([]interfaces.LayerInterface, 0),
		connectionstatus: make([]common.LayerStatus, 0),
		recvch:           make(chan *bufferpool.Buffer, 5),
		statusch:         make(chan common.LayerStatus, 1),
		status:           common.LayerDown,
	}

	go func(cm *Connmgr) {
		//
		// Close down the structure here
		//
		defer func() {
			cm.sendStatus(common.LayerClosed)
			close(cm.statusch)
			close(cm.recvch)
			for i := 0; i < len(cm.connections); i++ {
				cm.connections[i].Reset()
			}
		}()

		for {
			connlen := len(cm.connections)
			if connlen == 0 && cm.status != common.LayerDown {
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
				cm.connectionstatus = append(cm.connectionstatus, common.LayerDown)
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
						cm.connectionstatus[i-(2+l)] = recv.Interface().(common.LayerStatus)
						cm.updateStatus()
					}
				}
			}
		}
	}(cm.(*Connmgr))

	return cm, nil
}

// Send()
func (cm *Connmgr) Send(b *bufferpool.Buffer) {
	for i := 0; i < len(cm.connectionstatus); i++ {
		if cm.connectionstatus[i] == common.LayerUp {
			cm.connections[i].Send(b)
			return
		}
	}

	log.Warnf("Connmgr[%s] Send Drop...", cm.serial.Next())

}

// RecvCh()
func (cm *Connmgr) RecvCh() chan *bufferpool.Buffer {
	return cm.recvch
}

// StatusCh()
func (cm *Connmgr) StatusCh() chan common.LayerStatus {
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
func (cm *Connmgr) sendStatus(s common.LayerStatus) {
	cm.statusch <- s
}

// updateStatus()
// If one connection change to up up, send up, otherwise change to down
func (cm *Connmgr) updateStatus() {
	status := common.LayerDown

	for i := 0; i < len(cm.connectionstatus); i++ {
		if cm.connectionstatus[i] == common.LayerUp {
			status = common.LayerUp
			break
		}
	}
	if cm.status != status {
		cm.status = status
		cm.sendStatus(status)
	}
}
