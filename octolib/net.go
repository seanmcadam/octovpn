package octolib

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrBadIPPort = errors.New("Bad IP Port String")
var ErrBadNetMask = errors.New("Bad NetMask")
var ErrBadNetBits = errors.New("Bad NetBits")

var netbits = make(map[string]string)

func init() {
	netbits["255.255.255.255"] = "32"
	netbits["255.255.255.254"] = "31"
	netbits["255.255.255.252"] = "30"
	netbits["255.255.255.248"] = "29"
	netbits["255.255.255.240"] = "28"
	netbits["255.255.255.224"] = "27"
	netbits["255.255.255.198"] = "26"
	netbits["255.255.255.128"] = "25"
	netbits["255.255.255.0"] = "24"
	netbits["255.255.254.0"] = "23"
	netbits["255.255.252.0"] = "22"
	netbits["255.255.248.0"] = "21"
	netbits["255.255.240.0"] = "20"
	netbits["255.255.224.0"] = "19"
	netbits["255.255.198.0"] = "18"
	netbits["255.255.128.0"] = "17"
	netbits["255.255.0.0"] = "16"
}

func IP4netmask2netbits(netmask string) (string, error) {

	u, ok := netbits[netmask]
	if !ok {
		return "", ErrBadNetMask
	}

	return u, nil
}

func SplitAddr(ipport string) (ip string, p uint16, e error) {

	s := strings.Split(ipport, ":")
	if len(s) < 2 {
		e = errors.New(fmt.Sprintf("Bad IP Port String:%s", ipport))
	} else {
		var pint int
		ip = s[0]
		pint, e = strconv.Atoi(s[1])
		p = uint16(pint)
	}
	return
}
