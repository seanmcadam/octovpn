package settings

type ConfigProtoType string
type ConfigPortType uint16
type ConfigEndpointNameType string
type ConfigMtuType uint16
type ConfigWidthType uint8

const Width32 ConfigWidthType = 32
const Width64 ConfigWidthType = 64
const WidthDefault ConfigWidthType = Width64
const TCP ConfigProtoType = "tcp"
const TCP4 ConfigProtoType = "tcp4"
const TCP6 ConfigProtoType = "tcp6"
const UDP ConfigProtoType = "udp"
const UDP4 ConfigProtoType = "udp4"
const UDP6 ConfigProtoType = "udp6"
const ConfigMtuMin ConfigMtuType = 800
const ConfigMtuMax ConfigMtuType = 1500
const ConfigPortMin ConfigPortType = 1024
const ConfigPortMax ConfigPortType = 65534

type ConnectionStruct struct {
	Name   string          `json:"name"`
	Width  ConfigWidthType `json:"width,omitempty"`
	Proto  ConfigProtoType `json:"proto"`
	Host   string          `json:"host"`
	Port   ConfigPortType  `json:"port"`
	Mtu    ConfigMtuType   `json:"mtu,omitempty"`
	Auth   string          `json:"auth,omitempty"`
}

type ConfigSiteStruct struct {
	Sitename string             `json:"sitename"`
	Width    ConfigWidthType    `json:"width,omitempty"`
	Servers  []ConnectionStruct `json:"servers,omitempty"`
	Clients  []ConnectionStruct `json:"clients,omitempty"`
	Mtu      ConfigMtuType      `json:"mtu,omitempty"`
}

//
// This Node
// What sites am I connected with
//
type ConfigStruct struct {
	NodeName string             `json:"nodename"`
	Sites    []ConfigSiteStruct `json:"sites"`
}

func (c ConfigProtoType) String() string {
	return string(c)
}

func (c ConfigEndpointNameType) String() string {
	return string(c)
}
