package netlib

import (
	"fmt"
	"math/rand"
	"net"
	"testing"
)

var goodhostname = "google.com"
var badhostname = "_.google.com"

func TestValidIP(t *testing.T) {

	bogusip := "xxxxx"
	ipv4 := generateRandomIPv4()
	ipv6 := generateRandomIPv6()

	if ValidIP(bogusip) {
		t.Fatalf("Random IPv4: %s", bogusip)
	}

	if !ValidIP(ipv4) {
		t.Fatalf("Random IPv4: %s", ipv4)
	}

	if !ValidIP(ipv6) {
		t.Fatalf("Random IPv6: %s", ipv6)
	}

}

func TestValidIP4(t *testing.T) {

	bogusip := "yyyyy"
	ipv4 := generateRandomIPv4()
	ipv6 := generateRandomIPv6()

	if ValidIP4(bogusip) {
		t.Fatalf("Random IPv4: %s", bogusip)
	}

	if !ValidIP4(ipv4) {
		t.Fatalf("Random IPv4: %s", ipv4)
	}

	if ValidIP4(ipv6) {
		t.Fatalf("Random IPv4: %s", ipv6)
	}
}

func TestValidIP6(t *testing.T) {

	bogusip := "zzzzz"
	ipv4 := generateRandomIPv4()
	ipv6 := generateRandomIPv6()

	if ValidIP6(bogusip) {
		t.Fatalf("Random IPv4: %s", bogusip)
	}

	if ValidIP6(ipv4) {
		t.Fatalf("Random IPv6: %s", ipv4)
	}

	if !ValidIP6(ipv6) {
		t.Fatalf("Random IPv6: %s", ipv6)
	}
}

func TestValidIPHost(t *testing.T) {

	if !ValidIPHost(goodhostname) {
		t.Fatalf("IP goodhostname: %s", goodhostname)
	}
	if ValidIPHost(badhostname) {
		t.Fatalf("IP badhostname: %s", badhostname)
	}
}

func TestValidIP4Host(t *testing.T) {

	if !ValidIP4Host(goodhostname) {
		t.Fatalf("IP4 goodhostname: %s", goodhostname)
	}
	if ValidIP4Host(badhostname) {
		t.Fatalf("IP4 badhostname: %s", badhostname)
	}
}

func TestValidIP6Host(t *testing.T) {

	if !ValidIP6Host(goodhostname) {
		t.Fatalf("IP6 goodhostname: %s", goodhostname)
	}
	if ValidIP6Host(badhostname) {
		t.Fatalf("IP6 badhostname: %s", badhostname)
	}
}

func generateRandomIPv4() string {
	return fmt.Sprintf("%d.%d.%d.%d", rand.Intn(256), rand.Intn(256), rand.Intn(256), rand.Intn(256))
}

func generateRandomIPv6() string {
	ipv6 := make([]byte, 16)
	for i := 0; i < 16; i++ {
		ipv6[i] = byte(rand.Intn(256))
	}
	return net.IP(ipv6).String()
}
