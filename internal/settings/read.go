package settings

import (
	"encoding/json"
	"io"
	"os"

	"github.com/seanmcadam/octovpn/octolib/log"
)

func ReadConfig(filepath string) (config *ConfigStruct, err error) {
	// Open the JSON file

	file, err := os.Open(filepath)
	if err != nil {
		return nil, log.Errf("READ ERROR: opening file: %s", err)
	}
	defer file.Close()

	// Read the file into a byte slice
	rawFileByteValue, err := io.ReadAll(file)
	if err != nil {
		return nil, log.Errf("READ ERROR: reading file: %s", err)
	}

	// Create a variable of type Person
	var Config ConfigStruct

	// Unmarshal the byte slice into the person variable
	if err = json.Unmarshal(rawFileByteValue, &Config); err != nil {
		return nil, log.Errf("READ ERROR: marshaling JSON: %s", err)
	}

	if err = validateConfig(&Config); err != nil {
		return nil, err
	}

	return &Config, nil
}
