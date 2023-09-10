package settings

type ConfigProtoType string
type ConfigPortType uint16
type ConfigEndpointNameType string
type ConfigMtuType uint16

const TCP ConfigProtoType = "tcp"
const TCP4 ConfigProtoType = "tcp4"
const TCP6 ConfigProtoType = "tcp6"
const UDP ConfigProtoType = "udp"
const UDP4 ConfigProtoType = "udp4"
const UDP6 ConfigProtoType = "udp6"
const ConfigMtuMin ConfigMtuType = 1000
const ConfigMtuMax ConfigMtuType = 1500
const ConfigPortMin ConfigPortType = 1024
const ConfigPortMax ConfigPortType = 65534

type ConnectionStruct struct {
	Name  string          `json:"name"`
	Proto ConfigProtoType `json:"proto"`
	Host  string          `json:"host"`
	Port  ConfigPortType  `json:"port"`
	Mtu   ConfigMtuType   `json:"mtu,omitempty"`
	Auth  string          `json:"auth,omitempty"`
}

type ConfigSiteStruct struct {
	Sitename string             `json:"sitename"`
	Servers  []ConnectionStruct `json:"servers,omitempty"`
	Clients  []ConnectionStruct `json:"clients,omitempty"`
	Mtu      ConfigMtuType      `json:"mtu,omitempty"`
}

type ConfigStruct struct {
	NodeName string             `json:"nodename"`
	Sites    []ConfigSiteStruct `json:"sites"`
}
