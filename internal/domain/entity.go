package domain

import (
	"fmt"
	"time"
)

const (
	UDP = iota
	TCP
	TLS
)

// Phone configuration
// For examle used Yealink params for autoprovisioning
type PhoneConfig struct {
	model       string
	vendor      string
	macAddress  string
	ipAddress   string
	sipServer   string // account.1.sip_server.1.address
	sipPort     int
	codecs      []string // Потом надо будет добавить и прописать поддерживаемые.
	transport   int      // iota
	extensionID int      // Указывает на таблицу с номерами телефонов паролями и именами пользователей
	// account.1.user_name, account.1.password,account.1.auth_name
	lineState           bool // account.1.enable = 1 F.E.
	comliteState        bool
	updatedState        bool
	trustCtrl           bool
	autoProvisionState  bool   // Надо понять нужен ли этот параметр.
	autoProvisionServer string // Сам сервер. Надо понять нужен ли этот параметр.
	createdAt           *time.Time
	updatedAt           *time.Time
}

func (pc *PhoneConfig) SetVendor(val string) (*PhoneConfig, error) {
	pc.vendor = val
	return pc, nil
}
func (pc *PhoneConfig) SetModel(val string) (*PhoneConfig, error) {
	pc.model = val
	return pc, nil
}
func (pc *PhoneConfig) SetMACAddress(val string) (*PhoneConfig, error) {
	return pc, nil
}
func (pc *PhoneConfig) SetIPAddress(val string) (*PhoneConfig, error) {
	ip, err := validateIP(val)
	if err != nil {
		return pc, err
	}
	pc.ipAddress = ip
	return pc, nil
}
func (pc *PhoneConfig) SetSIPServer(val string) (*PhoneConfig, error) {
	server, err := validateIP(val)
	if err == nil {
		// Ip Валидный пишем сервер
		pc.sipServer = server
		return pc, nil
	}
	if validateDN(val) {
		pc.sipServer = val
		return pc, nil

	}
	return pc, fmt.Errorf("server address '%s' is neither a valid IP nor a valid domain name (IP error: %w)", val, err)
}
func (pc *PhoneConfig) SetSIPPort(val int) (*PhoneConfig, error) {
	if val < 0 || val > 65536 {
		return pc, fmt.Errorf("Invalid port number")
	}
	pc.sipPort = val
	return pc, nil
}
func (pc *PhoneConfig) SetTransport(val int) (*PhoneConfig, error) {
	if val < 0 || val > 3 {
		return pc, fmt.Errorf("Invalid transport mode")
	}
	pc.transport = val
	return pc, nil
}
func (pc *PhoneConfig) SetExtensionID(val int) (*PhoneConfig, error) {
	return pc, nil
}
func (pc *PhoneConfig) SetEnableAccoun(val bool) (*PhoneConfig, error) {
	return pc, nil
}
func (pc *PhoneConfig) SetChanged() (*PhoneConfig, error) {
	return pc, nil
}

func (pc *PhoneConfig) SetCodecs(val ...string) (*PhoneConfig, error) {
	return pc, nil
}
func (pc *PhoneConfig) SetProvisionServer(val string) (*PhoneConfig, error) {
	return pc, nil
}
func (pc *PhoneConfig) SetProvisioningEnable(val bool) (*PhoneConfig, error) {
	return pc, nil
}
func (pc *PhoneConfig) Complite() (bool, error) {
	return false, nil
}
