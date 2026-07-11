package domain

type PfoneSettings struct {
	id              int
	macAddress      string
	model           string
	vendor          string
	Lines           []LineSettings
	General         GeneralSettings
	NetworkSettings NetworkSettings
}

// Getters
func (p *PfoneSettings) LinesCount() int {
	return len(p.Lines)
}
func (p *PfoneSettings) MACAddress() string {
	return p.macAddress
}
func (p *PfoneSettings) Model() string {
	return p.model
}
func (p *PfoneSettings) Vendor() string {
	return p.vendor
}
func (p *PfoneSettings) GetLineSettings(index int) (*LineSettings, bool) {
	if index < 0 || index >= len(p.Lines) {
		return nil, false
	}
	line := p.Lines[index]
	return &line, true
}

// Setters
func (p *PfoneSettings) SetMACAddress(macAddress string) {
	if validateMAC(macAddress) {
		p.macAddress = macAddress
	}
}
func (p *PfoneSettings) SetModel(model string) {
	p.model = model
}
func (p *PfoneSettings) SetVendor(vendor string) {
	p.vendor = vendor
}
func (p *PfoneSettings) LinkLine(line LineSettings, slot int) error {
	if slot < 0 {
		return ErrInvalidSlot
	}
	if slot > len(p.Lines) {
		if p.Lines[slot] == line {
			p.Lines = append(p.Lines, line)
			return nil
		}
	}
	line.SetSlot(slot)
	p.Lines[slot] = line
	return nil
}
