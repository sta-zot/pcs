package config

import (
	"github.com/pelletier/go-toml"
)

// var configFile string = "/etc/pcs/config"
var configFile string = "S:\\WorkPlace\\pcs\\etc\\config"

type sipModel struct {
	Server    string `toml:"sip_server"`
	Port      int    `toml:"port"`
	Transport string `toml:"transport"`
	DtmfMode  string `toml:"dtmf_mode"`
	Expire    int    `toml:"expire"`
}

func (s *sipModel) GetServer() string {
	return s.Server
}
func (s *sipModel) GetPort() int {
	return s.Port
}
func (s *sipModel) GetTransport() string {
	return s.Transport
}
func (s *sipModel) GetDtmfMode() string {
	return s.DtmfMode
}
func (s *sipModel) GetExpire() int {
	return s.Expire
}

type generalModel struct {
	Timezone           string `toml:"timezone"`
	TimeFormat         int    `toml:"time_format"`
	Language           string `toml:"language"`
	WebPassword        string `toml:"web_password"`
	PhonePassword      string `toml:"phone_password"`
	ProvisionOnStart   bool   `toml:"provision_on_start"`
	EnableProvisioning bool   `toml:"enable_provisioning"`
	ProvisioningServer string `toml:"provisioning_server"`
}

func (g *generalModel) GetTimezone() string {
	return g.Timezone
}
func (g *generalModel) GetTimeFormat() int {
	return g.TimeFormat
}
func (g *generalModel) GetLanguage() string {
	return g.Language
}
func (g *generalModel) GetWebPassword() string {
	return g.WebPassword
}
func (g *generalModel) GetPhonePassword() string {
	return g.PhonePassword
}
func (g *generalModel) GetProvisionOnStart() bool {
	return g.ProvisionOnStart
}
func (g *generalModel) GetEnableProvisioning() bool {
	return g.EnableProvisioning
}
func (g *generalModel) GetProvisioningServer() string {
	return g.ProvisioningServer
}

type dataBaseConfModel struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	Database string `toml:"database"`
}

func (d *dataBaseConfModel) GetHost() string {
	return d.Host
}
func (d *dataBaseConfModel) GetPort() int {
	return d.Port
}
func (d *dataBaseConfModel) GetUser() string {
	return d.User
}
func (d *dataBaseConfModel) GetPassword() string {
	return d.Password
}
func (d *dataBaseConfModel) GetDatabase() string {
	return d.Database
}

type Config struct {
	Sip      sipModel          `toml:"sip"`
	General  generalModel      `toml:"general"`
	DataBase dataBaseConfModel `toml:"data_base"`
	Log      logModel          `toml:"log"`
}

func (c *Config) GetSip() sipModel {
	return c.Sip
}
func (c *Config) GetGeneral() generalModel {
	return c.General
}
func (c *Config) GetDataBase() dataBaseConfModel {
	return c.DataBase
}

func loadConfig(path string) (*Config, error) {
	var config Config
	data, err := toml.LoadFile(path)
	if err != nil {
		return nil, err
	}
	err = data.Unmarshal(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func NewConfig() (*Config, error) {
	config, err := loadConfig(configFile)
	if err != nil {
		return nil, err
	}
	return config, nil
}

type logModel struct {
	Level string `toml:"level"`
	Path  string `toml:"path"`
}
