package settings

import (
	"github.com/seanmcadam/loggy"
)

//
// Configs
//

func validateConfig(config *ConfigStruct) (err error) {

	if config == nil {
		return loggy.Err("NOT VALID: No Config Struct")
	}

	if len(config.NodeName) == 0 {
		return loggy.Err("NOT VALID: No NodeName defined")
	}

	if len(config.Sites) == 0 {
		return loggy.Err("NOT VALID: No Sites defined")
	}

	for i := 0; i < len(config.Sites); i++ {
		if err = validateSite(&config.Sites[i]); err != nil {
			loggy.Errorf("Site Error[%d]%s", i, err)
			return err
		}
	}

	return nil
}

func validateSite(site *ConfigSiteStruct) (err error) {
	if len(site.Sitename) == 0 {
		return loggy.Err("Not Site Name")
	}

	if site.Width != 0 {
		if site.Width != 32 && site.Width != 64 {
			return loggy.Errf("NOT VALID: Site Width is 63 or 32, Not:%d", int(site.Width))
		}
	}

	if len(site.Servers) == 0 && len(site.Clients) == 0 {
		return loggy.Err("Need One Server or Client specified")
	}

	if len(site.Servers) > 0 {
		for i := 0; i < len(site.Servers); i++ {
			if err = validateConnection(&site.Servers[i]); err != nil {
				return err
			}
		}
	}

	if len(site.Clients) > 0 {
		for i := 0; i < len(site.Clients); i++ {
			if err = validateConnection(&site.Clients[i]); err != nil {
				return err
			}
		}
	}

	if site.Mtu > 0 {
		if err = validateMtu(&site.Mtu); err != nil {
			return err
		}
	}

	return nil
}

func validateConnection(connection *ConnectionStruct) (err error) {

	if len(connection.Name) == 0 {
		return loggy.Err("Server Name Required")
	}

	if connection.Width != 0 {
		if connection.Width != 32 && connection.Width != 64 {
			return loggy.Errf("NOT VALID: Connection Width is 63 or 32, Not:%d", int(connection.Width))
		}
	}

	if len(connection.Host) == 0 {
		return loggy.Err("Host Name Required")
	}

	if err = validateProto(&connection.Proto); err != nil {
		return err
	}

	if err = validatePort(&connection.Port); err != nil {
		return err
	}

	if connection.Mtu > 0 {
		if err = validateMtu(&connection.Mtu); err != nil {
			return err
		}
	}

	return nil
}

func validateMtu(mtu *ConfigMtuType) (err error) {

	if *mtu < ConfigMtuMin || *mtu > ConfigMtuMax {
		return loggy.Errf("MTU Size:%d outside the MIN/MAX boundries %d/%d", *mtu, ConfigMtuMin, ConfigMtuMax)
	}

	return nil
}

func validatePort(port *ConfigPortType) (err error) {

	if *port < ConfigPortMin || *port > ConfigPortMax {
		return loggy.Errf("MTU Size:%d outside the MIN/MAX boundries %d/%d", *port, ConfigPortMin, ConfigPortMax)
	}

	return nil
}

func validateProto(proto *ConfigProtoType) (err error) {
	if string(*proto) == "" {
		return loggy.Errf("Protocol Required")
	}

	switch *proto {
	case TCP:
	case TCP4:
	case TCP6:
	case UDP:
	case UDP4:
	case UDP6:
	default:
		return loggy.Errf("Unknown Protocol:%s", *proto)
	}

	return nil
}
