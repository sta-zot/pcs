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
