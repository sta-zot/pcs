package domain

import "time"

type NetworkSettings struct {
	mode         NetworkMode // DHCP-0 / Static-1
	ip           string      // IP адрес
	subnetMask   string      // Маска подсети
	gateway      string      // Гейтвей
	primaryDNS   string      // Первичный DNS
	secondaryDNS string      // Вторичный DNS
	vlan         int         // VLAN по умолчанию 1 - VLAN выключен. 
	updatedAt    time.Time   // время последнего обновления
	hasChanged   bool        // для точечного обновления
}
