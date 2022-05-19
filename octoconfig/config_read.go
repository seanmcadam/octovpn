package octoconfig

import (
	"fmt"
	"strings"
)

func ReadConfigs() (configs ConfigV1) {
	configFile, err := ConfigGetVal(ConfigFilePath)
	if err != nil {
		panic(fmt.Sprintf("Getting Config File val failed: %s\n", err))
	}

	if configFile != "" { // if a config file was specified use it instead of flag params
		//
		// Config Configs
		//
		configjson, err := LoadConfiguration("config")
		if err != nil {
			panic(fmt.Sprintf("Failed to load config configuration: %v", err))
		}

		//var I ConfigInterface
		//var C ConfigConnection

		configs.Iface = parseInterface(configjson.Iface)
		configs.Conn = make(map[string]*ConfigConnection)

		for i, c := range configjson.Connections.(map[string]interface{}) {
			configs.Conn[i] = parseConnection(c)
		}
	}

	return configs
}

//
//
//
func parseInterface(data interface{}) (i *ConfigInterface) {

	i = &ConfigInterface{
		Name:    "",
		TunTap:  "",
		IP:      "",
		Netmask: "",
		MTU:     1400,
	}

	d, ok := data.(map[string]interface{})
	if !ok {
		panic("")
	}
	i.Name, ok = d["name"].(string)
	if !ok {
		panic("")
	}
	i.TunTap, ok = d["tuntap"].(string)
	if !ok {
		panic("")
	}
	i.IP, ok = d["ip"].(string)
	if !ok {
		panic("")
	}
	i.Netmask, ok = d["netmask"].(string)
	if !ok {
		panic("")
	}
	_, ok = d["mtu"]
	if ok {
		switch d["mtu"].(type) {
		case int:
			i.MTU = d["mtu"].(int)
		case float32:
			i.MTU = int(d["mtu"].(float32))
		case float64:
			i.MTU = int(d["mtu"].(float64))
		default:
			panic("")
		}
	}

	return i
}

//
//
//
func parseConnection(data interface{}) (c *ConfigConnection) {

	c = &ConfigConnection{
		Protocol: TCP,
		Hostname: "",
		Port:     0,
		MTU:      0,
	}

	d, ok := data.(map[string]interface{})
	if !ok {
		panic("")
	}
	protocol, ok := d["protocol"].(string)
	if !ok {
		panic("")
	}
	protocol = strings.ToLower(protocol)
	switch protocol {
	case "tcp":
		c.Protocol = TCP
	case "udp":
		c.Protocol = UDP
	default:
		panic("")
	}
	c.Hostname, ok = d["hostname"].(string)
	if !ok {
		panic("")
	}
	_, ok = d["port"]
	if !ok {
		panic("")
	}
	switch d["port"].(type) {
	case int:
		c.Port = d["port"].(int)
	case float32:
		c.Port = int(d["port"].(float32))
	case float64:
		c.Port = int(d["port"].(float64))
	default:
		panic("")
	}
	_, ok = d["mtu"]
	if ok {
		switch d["mtu"].(type) {
		case int:
			c.MTU = d["mtu"].(int)
		case float32:
			c.MTU = int(d["mtu"].(float32))
		case float64:
			c.MTU = int(d["mtu"].(float64))
		default:
			panic("")
		}
	}

	return c
}
