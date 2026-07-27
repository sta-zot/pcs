package domain

type DeviceInfo struct {
	vendor string
	model  string
	mac    string
	ip     string
}

// setters
func (d *DeviceInfo) SetModel(model string) {
	d.model = model
}
func (d *DeviceInfo) SetMAC(mac string) {
	d.mac = mac
}
func (d *DeviceInfo) SetVendor(vendor string) {
	d.vendor = vendor
}
func (d *DeviceInfo) SetIP(ip string) {
	d.ip = ip
}

// getters
func (d *DeviceInfo) Model() string {
	return d.model
}
func (d *DeviceInfo) MAC() string {
	return d.mac
}
func (d *DeviceInfo) Vendor() string {
	return d.vendor
}
func (d *DeviceInfo) IP() string {
	return d.ip
}

func NewDeviceInfo(model, vendor, mac, ip string) *DeviceInfo {
	return &DeviceInfo{
		model:  model,
		vendor: vendor,
		mac:    mac,
		ip:     ip,
	}
}
