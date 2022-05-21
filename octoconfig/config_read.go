package octoconfig

import (
	"strings"

	"github.com/seanmcadam/octovpn/octolib"
)

func ReadConfigs() (configs ConfigV1, e error) {
	configFile, err := ConfigGetVal(ConfigFilePath)
	if err != nil {
		return configs, octolib.ErrorLocationf("Getting Config File val failed: %s", err)
	}

	if configFile != "" { // if a config file was specified use it instead of flag params
		//
		// Config Configs
		//
		configjson, err := LoadConfiguration()
		if err != nil {
			return configs, octolib.ErrorLocationf("Failed to load config configuration: %v", err)
		}

		// Convert Configs

		configs.Iface, e = parseInterface(configjson.Iface)
		if e != nil {
			return configs, e
		}
		configs.Conn = make(map[string]*ConfigConnection)
		configs.List = make(map[string]*ConfigServer)

		if configjson.Connections != nil {
			for i, c := range configjson.Connections.(map[string]interface{}) {
				configs.Conn[i], e = parseConnection(c)
				if e != nil {
					return configs, e
				}
			}
		}

		if configjson.Listen != nil {
			for i, s := range configjson.Listen.(map[string]interface{}) {
				configs.List[i], e = parseServer(s)
				if e != nil {
					return configs, e
				}
			}
		}
	}

	return configs, e
}

//
//
//
func parseInterface(data interface{}) (i *ConfigInterface, e error) {

	var ok bool

	i = &ConfigInterface{
		Name:    "",
		TunTap:  "",
		IP:      "",
		Netmask: "",
		MTU:     1400,
	}

	d, ok := data.(map[string]interface{})
	if !ok {
		return nil, octolib.ErrorLocationf("wrong data type:%t", data)
	}

	i.Name, ok = d[string(configName)].(string)
	if !ok {
		return nil, octolib.ErrorLocationf("Cannot read Interface name")
	}

	i.TunTap, ok = d[string(configTunTap)].(string)
	if !ok {
		return nil, octolib.ErrorLocationf("Cannot read Interface device type")
	}
	i.TunTap = strings.ToUpper(i.TunTap)
	switch i.TunTap {
	case "TUN":
	case "TAP":
	default:
		return nil, octolib.ErrorLocationf("Bad Interface device type:%t", i.TunTap)
	}

	i.IP, ok = d[string(configIP)].(string)
	if !ok {
		return nil, octolib.ErrorLocationf("Cannot read Interface IP")
	}

	i.Netmask, ok = d[string(configNetmask)].(string)
	if !ok {
		return nil, octolib.ErrorLocationf("Cannot read Interface Netmask")
	}

	_, ok = d[string(configMTU)]
	if ok {
		switch d[string(configMTU)].(type) {
		case int:
			i.MTU = d[string(configMTU)].(int)
		case float32:
			i.MTU = int(d[string(configMTU)].(float32))
		case float64:
			i.MTU = int(d[string(configMTU)].(float64))
		default:
			return nil, octolib.ErrorLocationf("Bad MTU type:%t", d[string(configMTU)])
		}
	}

	if i.MTU < 512 || i.MTU > 1500 {
		return nil, octolib.ErrorLocationf("Bad MTU range:%d", i.MTU)
	}

	return i, e
}

//
//
//
func parseConnection(data interface{}) (c *ConfigConnection, e error) {

	c = &ConfigConnection{
		Protocol: TCP,
		Hostname: "",
		Port:     0,
		MTU:      0,
	}

	d, ok := data.(map[string]interface{})
	if !ok {
		return nil, octolib.ErrorLocationf("Cannot create map")
	}

	protocol, ok := d[string(configProtocol)].(string)
	if !ok {
		return nil, octolib.ErrorLocationf("Cannot read Protocol")
	}
	protocol = strings.ToLower(protocol)
	switch protocol {
	case "tcp":
		c.Protocol = TCP
	case "udp":
		c.Protocol = UDP
	default:
		return nil, octolib.ErrorLocationf("Bad Protocol:%s", protocol)
	}
	c.Hostname, ok = d[string(configHostname)].(string)
	if !ok {
		return nil, octolib.ErrorLocationf("Cannot read Hostname")
	}
	if c.Hostname == "" {
		return nil, octolib.ErrorLocationf("Empty Hostname")
	}
	_, ok = d[string(configPort)]
	if !ok {
		return nil, octolib.ErrorLocationf("Cannot read Port")
	}
	switch d[string(configPort)].(type) {
	case int:
		c.Port = d[string(configPort)].(int)
	case float32:
		c.Port = int(d[string(configPort)].(float32))
	case float64:
		c.Port = int(d[string(configPort)].(float64))
	default:
		return nil, octolib.ErrorLocationf("Bad Port type:%t", d[string(configPort)])
	}
	_, ok = d[string(configMTU)]
	if ok {
		switch d[string(configMTU)].(type) {
		case int:
			c.MTU = d[string(configMTU)].(int)
		case float32:
			c.MTU = int(d[string(configMTU)].(float32))
		case float64:
			c.MTU = int(d[string(configMTU)].(float64))
		default:
			return nil, octolib.ErrorLocationf("Bad MTU type:%t", d[string(configMTU)])
		}
	}

	return c, e
}

//
//
//
func parseServer(data interface{}) (s *ConfigServer, e error) {

	s = &ConfigServer{
		Protocol: TCP,
		IP:       "",
		Port:     0,
		MTU:      0,
	}

	d, ok := data.(map[string]interface{})
	if !ok {
		return nil, octolib.ErrorLocationf("Cannot create map")
	}

	protocol, ok := d[string(configProtocol)].(string)
	if !ok {
		return nil, octolib.ErrorLocationf("")
	}
	protocol = strings.ToLower(protocol)
	switch protocol {
	case "tcp":
		s.Protocol = TCP
	case "udp":
		s.Protocol = UDP
	default:
		return nil, octolib.ErrorLocationf("")
	}
	s.IP, ok = d[string(configIP)].(string)
	if !ok {
		return nil, octolib.ErrorLocationf("")
	}
	_, ok = d[string(configPort)]
	if !ok {
		return nil, octolib.ErrorLocationf("")
	}
	switch d[string(configPort)].(type) {
	case int:
		s.Port = d[string(configPort)].(int)
	case float32:
		s.Port = int(d[string(configPort)].(float32))
	case float64:
		s.Port = int(d[string(configPort)].(float64))
	default:
		return nil, octolib.ErrorLocationf("")
	}
	_, ok = d[string(configMTU)]
	if ok {
		switch d[string(configMTU)].(type) {
		case int:
			s.MTU = d[string(configMTU)].(int)
		case float32:
			s.MTU = int(d[string(configMTU)].(float32))
		case float64:
			s.MTU = int(d[string(configMTU)].(float64))
		default:
			return nil, octolib.ErrorLocationf("")
		}
	}

	return s, e
}
