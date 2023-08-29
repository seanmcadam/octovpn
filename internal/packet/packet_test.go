package packet

import (
	"fmt"
	"testing"
)

func TestPacket_strings(t *testing.T) {
	_ = fmt.Sprintf("%s", ROUTE_TYPE_RAW)
	_ = fmt.Sprintf("%s", ROUTE_TYPE_ROUTE)
	_ = fmt.Sprintf("%s", ROUTE_TYPE_ETH)
	_ = fmt.Sprintf("%s", ROUTE_TYPE_IP4)
	_ = fmt.Sprintf("%s", ROUTE_TYPE_IP6)
	_ = fmt.Sprintf("%s", SITE_TYPE_RAW)
	_ = fmt.Sprintf("%s", SITE_TYPE_AUTH)
	_ = fmt.Sprintf("%s", SITE_TYPE_PARENT)
	_ = fmt.Sprintf("%s", SITE_TYPE_ACK)
	_ = fmt.Sprintf("%s", SITE_TYPE_NAK)
	_ = fmt.Sprintf("%s", SITE_TYPE_PING64)
	_ = fmt.Sprintf("%s", SITE_TYPE_PONG64)
	_ = fmt.Sprintf("%s", SITE_TYPE_ERROR)
	_ = fmt.Sprintf("%s", CHAN_TYPE_RAW)
	_ = fmt.Sprintf("%s", CHAN_TYPE_AUTH)
	_ = fmt.Sprintf("%s", CHAN_TYPE_PARENT)
	_ = fmt.Sprintf("%s", CHAN_TYPE_ACK)
	_ = fmt.Sprintf("%s", CHAN_TYPE_NAK)
	_ = fmt.Sprintf("%s", CHAN_TYPE_PING64)
	_ = fmt.Sprintf("%s", CHAN_TYPE_PONG64)
	_ = fmt.Sprintf("%s", CHAN_TYPE_ERROR)
	_ = fmt.Sprintf("%s", CONN_TYPE_RAW)
	_ = fmt.Sprintf("%s", CONN_TYPE_AUTH)
	_ = fmt.Sprintf("%s", CONN_TYPE_PARENT)
	_ = fmt.Sprintf("%s", CONN_TYPE_ACK)
	_ = fmt.Sprintf("%s", CONN_TYPE_NAK)
	_ = fmt.Sprintf("%s", CONN_TYPE_PING64)
	_ = fmt.Sprintf("%s", CONN_TYPE_PONG64)
	_ = fmt.Sprintf("%s", CONN_TYPE_ERROR)
}
