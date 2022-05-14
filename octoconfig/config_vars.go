package octoconfig

//
// Define all configuration variables here
//

type ConfigConst string

const (
	ConfigHelp     ConfigConst = "ConfigHelp"
	ConfigFilePath ConfigConst = "ConfigFile"
	IFaceName      ConfigConst = "IFaceName"
	LogLevel       ConfigConst = "LogLevel"
)

// Config Structure
type configitem struct {
	//
	// Internal Name for the flag value
	//
	flagname string
	//
	// Description of how the flag is used
	//
	flagusage string
	//
	// The Name of the environmental variable (if applicable)
	//
	envname string
	//
	// The Name of the configuration file variable (if applicable)
	//
	configname string
	//
	// The default value of the variable (used to inialize the variable)
	//
	defval string
	//
	// Stored value of the configuration variable
	//
	configval *string
	//
	// Stored value of the environmental variable
	//
	envval *string
	//
	// Stored value of the command line variable
	//
	flagval *string
}

//
// Define all Configuration constants that can be read in from the command line and the environment
//
var ConfigMap = map[ConfigConst]configitem{
	ConfigHelp: {
		flagname:   "help",
		flagusage:  "Display The help Screen",
		envname:    "",
		configname: "",
		defval:     "",
		configval:  nil,
		envval:     nil,
		flagval:    nil,
	},
	ConfigFilePath: {
		flagname:   "configfile",
		flagusage:  "Configuration File Name, default, environment, or command line only",
		envname:    "ConfigFilePath",
		configname: "",
		defval:     "octovpn.json",
		configval:  nil,
		envval:     nil,
		flagval:    nil,
	},
	IFaceName: {
		flagname:   "interfacename",
		flagusage:  "VPN Interface Name",
		envname:    "IFaceName",
		configname: "",
		defval:     "octovpn",
		configval:  nil,
		envval:     nil,
		flagval:    nil,
	},
	LogLevel: {
		flagname:   "loglevel",
		flagusage:  "Logging Level: 0-7",
		envname:    "LogLevel",
		configname: "LogLevel",
		defval:     "5",
		configval:  nil,
		envval:     nil,
		flagval:    nil,
	},
}
