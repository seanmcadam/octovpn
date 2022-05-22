package octoconfig

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type DeviceType string
type ConnectionProtocol string
type ConfigMapField string

const TUN DeviceType = "TUN"
const TAP DeviceType = "TAP"

const TCP ConnectionProtocol = "tcp"
const UDP ConnectionProtocol = "udp"

const configHostname ConfigMapField = "hostname"
const configIP ConfigMapField = "ip"
const configMTU ConfigMapField = "mtu"
const configName ConfigMapField = "name"
const configNetmask ConfigMapField = "netmask"
const configPort ConfigMapField = "port"
const configProtocol ConfigMapField = "protocol"
const configTunTap ConfigMapField = "tuntap"

type TargMap map[string]*ConfigTarget
type ListMap map[string]*ConfigListen

type ConfigJson struct {
	Version int         `json:"version,omitempty"`
	Iface   interface{} `json:"interface"`
	Target  interface{} `json:"target,omitempty"`
	Listen  interface{} `json:"listen,omitempty"`
}

type ConfigV1 struct {
	Iface *ConfigInterface
	Targ  TargMap
	List  ListMap
}

type ConfigInterface struct {
	Name    string `json:"name"`
	TunTap  string `json:"tuntap"`
	IP      string `json:"ip"`
	Netmask string `json:"netmask"`
	MTU     int    `json:"mtu,omitempty"`
}

type ConfigTarget struct {
	Protocol ConnectionProtocol
	Hostname string
	Port     uint16
	MTU      uint16
}

type ConfigListen struct {
	Protocol ConnectionProtocol
	IP       string
	Port     uint16
	MTU      uint16
}

//
// Reads in the configuration map
// Inializes the environmental variables
//
func init() {

	//
	// Read in command line flags (config_vars)
	//
	for i, v := range ConfigMap {
		if v.flagname != "" {
			v.flagval = flag.String(v.flagname, "", v.flagusage)
			ConfigMap[i] = v
		}
	}

	flag.Parse()

	//
	// Read in environmental variables
	//
	for i, v := range ConfigMap {
		j := os.Getenv(v.envname)
		if j != "" {
			v.envval = &j
			ConfigMap[i] = v
		}
	}

}

func ConfigGetFlagName(cc ConfigConst) (v string, e error) {
	if _, ok := ConfigMap[cc]; ok {
		v = ConfigMap[cc].flagname
	} else {
		e = fmt.Errorf("undefined config constant: %s", cc)
	}
	return
}

func ConfigGetEnvName(cc ConfigConst) (v string, e error) {
	if _, ok := ConfigMap[cc]; ok {
		v = ConfigMap[cc].envname
	} else {
		e = fmt.Errorf("undefined config constant: %s", cc)
	}
	return
}

func ConfigGetConfigName(cc ConfigConst) (v string, e error) {
	if _, ok := ConfigMap[cc]; ok {
		v = ConfigMap[cc].configname
	} else {
		e = fmt.Errorf("undefined config constant: %s", cc)
	}
	return
}

//
// MustGet()
// looks for the config key and returns a string value no matter what.
//
func MustGet(cc ConfigConst) string {
	val, _ := ConfigGetVal(cc)

	return val
}

func ConfigGetVal(cc ConfigConst) (v string, e error) {
	if val, ok := ConfigMap[cc]; ok {
		if val.flagval != nil && *val.flagval != "" {
			v = *val.flagval
		} else if val.envval != nil {
			v = *val.envval
		} else if val.configval != nil {
			v = *val.configval
		} else {
			v = val.defval
		}
	} else {
		e = fmt.Errorf("undefined config constant: %s", cc)
	}

	return
}

//
// Local File takes precidence over remote download config
// func LoadConfiguration(pkg string) (data map[string]interface{}, err error) {
//
func LoadConfiguration() (data ConfigJson, err error) {
	currDir, _ := os.Getwd()
	defer os.Chdir(currDir)

	filePath, err := ConfigGetVal(ConfigFilePath)
	if err != nil {
		panic(fmt.Errorf("error retrieving config file variable: %s", err))
	}
	file := filepath.Base(filePath)
	filePath = filepath.Dir(filePath)
	os.Chdir(filePath)

	configFile, err := os.Open(file)
	if err != nil {
		return data, err
	}
	defer configFile.Close()
	byteValue, _ := ioutil.ReadAll(configFile)

	err = json.Unmarshal(byteValue, &data)
	return data, err
}

//
//
//
func DownloadConfig(fullURLFile string) {
	fileURL, err := url.Parse(fullURLFile)
	if err != nil {
		log.Fatal(err)
	}
	path := fileURL.Path
	segments := strings.Split(path, "/")
	fileName := segments[len(segments)-1]

	var file *os.File
	if flag.Lookup("test.v") == nil { // called from main.go
		file, err = os.Create("./configurationmanager/" + fileName)
		if err != nil {
			log.Fatal(err)
		}
	} else { // called from test srcipt
		file, err = os.Create(fileName)
		if err != nil {
			log.Fatal(err)
		}
	}

	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	resp, err := client.Get(fullURLFile)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	size, err := io.Copy(file, resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fmt.Printf("Downloaded a file %s with size %d\n", fileName, size)
}
