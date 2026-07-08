package domain

import (
	"time"
)

const (
	UDP = iota
	TCP
	TLS
)

type NetworkMode int

const (
	DHCP NetworkMode = iota
	Static
)

type PhoneConfig struct {
	network      NetworkSettings
	sip          SipSettings
	provisioning ProvisioningSettings
	general      GeneralSettings
	line         LineSettings
	id           int
	macAddr      string
	vendor       string
	model        string
	userAgent    string
	//Используются для выборки по состоянию конфигурации в БД
	isComplete bool
	hasChanged bool
	lastSeen   time.Time
	changedAt  time.Time
	updatedAt  time.Time
}

// ============================================
// ГЕТТЕРЫ (все приватные поля)
// ============================================

func (p *PhoneConfig) ID() int                            { return p.id }
func (p *PhoneConfig) MACAddr() string                    { return p.macAddr }
func (p *PhoneConfig) Vendor() string                     { return p.vendor }
func (p *PhoneConfig) Model() string                      { return p.model }
func (p *PhoneConfig) UserAgent() string                  { return p.userAgent }
func (p *PhoneConfig) HasChanged() bool                   { return p.hasChanged }
func (p *PhoneConfig) LastSeen() time.Time                { return p.lastSeen }
func (p *PhoneConfig) ChangedAt() time.Time               { return p.changedAt }
func (p *PhoneConfig) UpdatedAt() time.Time               { return p.updatedAt }
func (p *PhoneConfig) Network() NetworkSettings           { return p.network }
func (p *PhoneConfig) Sip() SipSettings                   { return p.sip }
func (p *PhoneConfig) Provisioning() ProvisioningSettings { return p.provisioning }
func (p *PhoneConfig) General() GeneralSettings           { return p.general }
func (p *PhoneConfig) Line() LineSettings                 { return p.line }

// ============================================
// СЕТТЕРЫ (только для изменяемых полей)
// ============================================

// ID устанавливается только при загрузке из БД
func (p *PhoneConfig) SetID(id int) {
	p.id = id
}

// MAC и Vendor не меняются после создания но может потребоваться перегонять из модели БД
func (p *PhoneConfig) SetMACAddr(mac string) {
	p.macAddr = mac
	p.touch()
}
func (p *PhoneConfig) SetVendor(vendor string) {
	p.vendor = vendor
	p.touch()
}

func (p *PhoneConfig) SetModel(model string) {
	p.model = model
	p.touch()
}

func (p *PhoneConfig) SetUserAgent(ua string) {
	p.userAgent = ua
	p.touch()
}

func (p *PhoneConfig) SetHasChanged(v bool) {
	p.hasChanged = v
	p.touch()
}

func (p *PhoneConfig) SetLastSeen(t time.Time) {
	p.lastSeen = t
}

// ============================================
// БИЗНЕС-МЕТОДЫ
// ============================================

// Вычислямое свойство для определения состояния конфигурации
// Возвращает true, если конфигурация телефона завершена и выставляет isComplete в TRUE
func (p *PhoneConfig) IsComplete() bool { return p.isComplete }

// MarkChanged помечает конфиг как изменённый и сбрасывает флаг после отправки
func (p *PhoneConfig) MarkChanged() {
	p.hasChanged = true
	p.changedAt = time.Now()
	p.touch()
}

func (p *PhoneConfig) MarkSynced() {
	p.hasChanged = false
	p.touch()
}

// touch обновляет updatedAt
func (p *PhoneConfig) touch() {
	p.updatedAt = time.Now()
}
