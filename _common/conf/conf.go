package conf

import (
	"encoding/json"
	"fmt"
	"os"
)

const confFile = "./dev.conf"

type DevConfig struct {
	Port          *string   `json:"port"`
	DeployPath    *string   `json:"deployPath"`
	JsPath        *string   `json:"jsPath"`
	CssPath       *string   `json:"cssPath"`
	EntryPoints   *[]string `json:"entryPoints"`
	Bundle        *bool     `json:"bundle"`
	WebComponents *struct {
		StylesDir *string `json: stylesDir`
	} `json:"webComponents`
}

func Get() *DevConfig {

	if _, err := os.Stat(confFile); err != nil && os.IsNotExist(err) {
		return nil
	} else if err != nil {
		fmt.Printf("\nConfig file: \x1b[31m%v: %v\x1b[0m", "Error opening the file", err)
		return nil
	}

	data, err := os.ReadFile(confFile)
	if err != nil {
		fmt.Printf("\nConfig file: \x1b[31m%v: %v\x1b[0m", "Error reading the file", err)
		return nil
	}

	var config DevConfig
	if err = json.Unmarshal(data, &config); err != nil {
		fmt.Printf("\nConfig file: \x1b[31m%v: %v\x1b[0m", "Unmarshalling error", err)
		return nil
	}

	return &config
}
