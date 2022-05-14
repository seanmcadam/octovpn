package config

import (
	"fmt"
)

type ConfigRead struct {
	Scheduler  string
	ConfigFile string
	IFaceName  string
}

func ReadConfigs() (configs ConfigRead) {
	configFile, err := ConfigGetVal(ConfigFilePath)
	if err != nil {
		panic(fmt.Sprintf("Getting Config File val failed: %s\n", err))
	}

	if configFile != "" { // if a config file was specified use it instead of flag params
		//
		// Config Configs
		//
		configConfig, err := LoadConfiguration("config")
		if err != nil {
			panic(fmt.Sprintf("Failed to load config configuration: %v", err))
		}

		configs.IFaceName = configConfig["IFaceName"].(string)

		//
		// Logging Configs
		//
		//		logLevelStr, err := ConfigGetVal(LogLevel)
		//		if err != nil {
		//			panic(fmt.Sprintf("Getting Log Level val failed: %s\n", err))
		//		}
		//		configs.LogLevel, err = strconv.Atoi(LevelStr)
		//		if err != nil {
		//			panic(fmt.Sprintf("Converting Log Level string to int failed: %s\n", err))
		//		}
		//		configs.LogFilePath, err = ConfigGetVal(LogFilePath)
		//		if err != nil {
		//			panic(fmt.Sprintf("Getting Log File Path val failed: %s\n", err))
		//		}
	}

	return configs
}
