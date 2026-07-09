package domain

import (
	"time"
)

type NetworkSettings struct {
	static       bool      // DHCP-0 / Static-1
	ipAddress    string    // IP адрес
	subnetMask   string    // Маска подсети
	gateway      string    // Гейтвей
	primaryDNS   string    // Первичный DNS
	secondaryDNS string    // Вторичный DNS
	vlanID       int       // VLAN по умолчанию 1 - VLAN выключен.
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
