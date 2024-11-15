package common

import (
	"encoding/json"
	"fmt"
	"os"
)

type devConfig struct {
	Port          *string   `json:"port"`
	DeployPath    *string   `json:"deployPath"`
	JsPath        *string   `json:"jsPath"`
	CssPath       *string   `json:"cssPath"`
	EntryPoints   *[]string `json:"entryPoints"`
	Bundle        *bool     `json:"bundle"`
	CreateGzip    *bool     `json:"createGzip"`
	WebComponents *struct {
		StylesDir *string `json:"stylesDir"`
	} `json:"webComponents"`
}

var Config *devConfig

func LoadConf() {
	const confFile = "./dev.conf"

	if _, err := os.Stat(confFile); err != nil && os.IsNotExist(err) {
		return
	} else if err != nil {
		fmt.Printf("\nConfig file: \x1b[31m%v: %v\x1b[0m", "Error opening the file", err)
		return
	}

	data, err := os.ReadFile(confFile)
	if err != nil {
		fmt.Printf("\nConfig file: \x1b[31m%v: %v\x1b[0m", "Error reading the file", err)
		return
	}

	var config devConfig
	if err = json.Unmarshal(data, &config); err != nil {
		fmt.Printf("\nConfig file: \x1b[31m%v: %v\x1b[0m", "Unmarshalling error", err)
		return
	}

	Config = &config
}
