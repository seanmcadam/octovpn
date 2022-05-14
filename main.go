package main

import (
	"fmt"

	"github.com/seanmcadam/octovpn/config"
)

func main() {

	config.ConfigGetVal(IFaceName)
	fmt.Printf("Running...\n")
}
