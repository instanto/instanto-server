package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config struct {
	Installed bool `json:"installed"`
	Port      int  `json:"port"`

	MediaDir string `json:"media_dir"`

	Secret string `json:"secret"`

	ServeWebApps              bool   `json:"serve_web_apps"`
	WebAppAdminDir            string `json:"web_app_admin_dir"`
	WebAppAdminURL            string `json:"web_app_admin_url"`
	WebAppPublicDir           string `json:"web_app_public_dir"`
	WebAppPublicURL           string `json:"web_app_public_url"`
	CORSEnabled               bool   `json:"cors_enabled"`
	AccessControlAllowOrigin  string `json:"access_control_allow_origin"`
	AccessControlAllowMethods string `json:"access_control_allow_methods"`
	AccessControlAllowHeaders string `json:"access_control_allow_headers"`

	DSN string `json:""`
}

type ConfigProvider struct {
	configFile string
}

func NewConfigProvider(configFile string) (*ConfigProvider, error) {
	return &ConfigProvider{configFile}, nil
}

func (cp *ConfigProvider) Parse() (*Config, error) {
	fd, err := os.Open(cp.configFile)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(fd)
	if err != nil {
		return nil, err
	}
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
