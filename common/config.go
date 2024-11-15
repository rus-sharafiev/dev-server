package common

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

var Config devConfig
