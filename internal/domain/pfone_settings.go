package domain

import (
	"errors"
	"fmt"
)

type PfoneSettings struct {
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
func (p *PfoneSettings) Complete() (bool, error) {
	err := p.checkComplete()
	return p.isComplete, err
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

// Helper functions
func (p *PfoneSettings) checkComplete() error {
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

func NewPfoneSettings(info DeviceInfo) *PfoneSettings {
	return &PfoneSettings{
		macAddress:      info.MAC(),
		model:           info.Model(),
		vendor:          info.Vendor(),
		NetworkSettings: *NewNetworkSettings(info.IP()),
	}
}
