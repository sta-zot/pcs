package domain

import (
	"errors"
	"time"
)

type NetworkSettings struct {
	static       bool      // DHCP-0 / Static-1
	ipAddress    string    // IP адрес
	subnetMask   string    // Маска подсети
	gateway      string    // Шлюз
	primaryDNS   string    // Первичный DNS
	secondaryDNS string    // Вторичный DNS
	vlanID       int       // VLAN по умолчанию 0 - означает что vlan  выключен, 1-4096 - VLAN включен.
	isComplete   bool      // Флаг, указывающий, что данные настроек заполнены полностью.
	updatedAt    time.Time // время последнего обновления
	hasChanged   bool      // для точечного обновления
}

// getters
func (ns *NetworkSettings) Static() bool         { return ns.static }
func (ns *NetworkSettings) IPAddress() string    { return ns.ipAddress }
func (ns *NetworkSettings) SubnetMask() string   { return ns.subnetMask }
func (ns *NetworkSettings) Gateway() string      { return ns.gateway }
func (ns *NetworkSettings) PrimaryDNS() string   { return ns.primaryDNS }
func (ns *NetworkSettings) SecondaryDNS() string { return ns.secondaryDNS }
func (ns *NetworkSettings) VlanID() int          { return ns.vlanID }
func (ns *NetworkSettings) UpdatedAt() time.Time { return ns.updatedAt }
func (ns *NetworkSettings) HasChanged() bool     { return ns.hasChanged }
func (ns *NetworkSettings) Complete() (bool, error) {
	err := ns.checkComplete()
	return ns.isComplete, err
}
func (ns *NetworkSettings) VLANEnabled() bool { return ns.vlanID != 0 }

// setters
func (ns *NetworkSettings) SetStatic(static bool) {
	ns.static = static
	ns.markChanged()
}
func (ns *NetworkSettings) SetIPAddress(ip string) {
	ip, err := validateIP(ip)
	if err != nil {
		return
	}
	ns.ipAddress = ip
	ns.markChanged()
}
func (ns *NetworkSettings) SetSubnetMask(mask string) {
	mask, err := validateIP(mask)
	if err != nil {
		return
	}
	ns.subnetMask = mask
	ns.markChanged()
}
func (ns *NetworkSettings) SetGateway(gateway string) {
	gateway, err := validateIP(gateway)
	if err != nil {
		return
	}
	ns.gateway = gateway
	ns.markChanged()
}
func (ns *NetworkSettings) SetPrimaryDNS(dns string) {
	dns, err := validateIP(dns)
	if err != nil {
		return
	}
	ns.primaryDNS = dns
	ns.markChanged()
}
func (ns *NetworkSettings) SetSecondaryDNS(dns string) {
	dns, err := validateIP(dns)
	if err != nil {
		return
	}
	ns.secondaryDNS = dns
	ns.markChanged()
}
func (ns *NetworkSettings) SetVlanID(vlan int) {
	if vlan < 1 || vlan > 4094 {
		return
	}
	ns.vlanID = vlan
	ns.markChanged()
}

// Helper function to mark the settings as changed and update the updatedAt timestamp
func (ns *NetworkSettings) touch() {
	ns.updatedAt = time.Now()
}
func (ns *NetworkSettings) markChanged() {
	ns.hasChanged = true
	ns.touch()
}
func (ns *NetworkSettings) MarkSynced() {
	ns.hasChanged = false

}

// Constructor for NetworkSettings
func NewNetworkSettings(ip string) *NetworkSettings {
	return &NetworkSettings{
		ipAddress:  ip,
		static:     false,
		updatedAt:  time.Now(),
		hasChanged: true,
	}
}

func (ns *NetworkSettings) checkComplete() error {
	if !ns.Static() {
		ns.isComplete = true
		return nil
	}
	if ns.ipAddress == "" {
		ns.isComplete = false
		return errors.New("ip address is required")
	}
	if ns.subnetMask == "" {
		ns.isComplete = false
		return errors.New("subnet mask is required")
	}
	if ns.gateway == "" {
		ns.isComplete = false
		return errors.New("gateway is required")
	}
	if ns.primaryDNS == "" {
		ns.isComplete = false
		return errors.New("primary DNS is required")
	}
	ns.isComplete = true
	return nil
}
