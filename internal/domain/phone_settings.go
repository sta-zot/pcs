package domain

import (
	"errors"
	"fmt"
)

type PhoneSettings struct {
	id              int
	macAddress      string
	model           string
	vendor          string
	Lines           []LineSettings
	General         GeneralSettings
	NetworkSettings NetworkSettings
	isComplete      bool
}

// Getters
func (p *PhoneSettings) LinesCount() int {
	return len(p.Lines)
}
func (p *PhoneSettings) MACAddress() string {
	return p.macAddress
}
func (p *PhoneSettings) Model() string {
	return p.model
}
func (p *PhoneSettings) Vendor() string {
	return p.vendor
}
func (p *PhoneSettings) GetLineSettings(index int) (*LineSettings, bool) {
	if index < 0 || index >= len(p.Lines) {
		return nil, false
	}
	line := p.Lines[index]
	return &line, true
}
func (p *PhoneSettings) Complete() (bool, error) {
	err := p.checkComplete()
	return p.isComplete, err
}

// Setters
func (p *PhoneSettings) SetMACAddress(macAddress string) {
	if validateMAC(macAddress) {
		p.macAddress = macAddress
	}
}
func (p *PhoneSettings) SetModel(model string) {
	p.model = model
}
func (p *PhoneSettings) SetVendor(vendor string) {
	p.vendor = vendor
}
func (p *PhoneSettings) LinkLine(line LineSettings, slot int) error {
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

// Helper functions
func (p *PhoneSettings) checkComplete() error {
	if p.macAddress == "" {
		p.isComplete = false
		return errors.New("missing required fields: MAC Address")
	}
	if p.model == "" {
		p.isComplete = false
		return errors.New("missing required fields: Model")
	}
	if p.vendor == "" {
		p.isComplete = false
		return errors.New("missing required fields: Vendor")
	}

	if len(p.Lines) == 0 {
		p.isComplete = false
		return errors.New("no lines attached")
	}
	for _, line := range p.Lines {
		if ok, err := line.Complete(); !ok || err != nil {
			p.isComplete = false
			return fmt.Errorf("incomplete line %d : %w", line.Slot(), err)
		}
	}
	if ok, err := p.NetworkSettings.Complete(); !ok || err != nil {
		p.isComplete = false
		return fmt.Errorf("network settings incomplete: %w", err)
	}

	p.isComplete = true
	return nil
}

func NewPhoneSettings(info DeviceInfo) *PhoneSettings {
	return &PhoneSettings{
		macAddress:      info.MAC(),
		model:           info.Model(),
		vendor:          info.Vendor(),
		NetworkSettings: *NewNetworkSettings(info.IP()),
	}
}
