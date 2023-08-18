package octotypes

import (
	"math/rand"
)

type NetworkPort uint16

func GetRandomNetworkPort() (n NetworkPort) {
	i := rand.Intn(50000) + 10000
	n = NetworkPort(i)
	return n
}
