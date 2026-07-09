package domain

import "time"

type GeneralSettings struct {
	timezone          int       // Часовой пояс
	timeFormat        string    // Формат времени (12h/24h)
	language          string    // Язык интерфейса
	webPassword       string    // Пароль для веб-интерфейса
	phonePassword     string    // Пароль для телефона
	provisionOnStart  bool      // static.auto_provision.power_on
	provisionURL      string    // static.auto_provision.server.url =
	provisionInterval int       // static.auto_provision.flexible.interval =
	provisionEnable   bool      //static.auto_provision.flexible.enable =
	hasChanged        bool      // для точечного обновления
	updatedAt         time.Time // время последнего обновления
}

// Getters
func (gs *GeneralSettings) Timezone() int {
	return gs.timezone
}

func (gs *GeneralSettings) TimeFormat() string {
	return gs.timeFormat
}

func (gs *GeneralSettings) Language() string {
	return gs.language
}

func (gs *GeneralSettings) WebPassword() string {
	return gs.webPassword
}

func (gs *GeneralSettings) PhonePassword() string {
	return gs.phonePassword
}

func (gs *GeneralSettings) ProvisionOnStart() bool {
	return gs.provisionOnStart
}

func (gs *GeneralSettings) ProvisionURL() string {
	return gs.provisionURL
}

func (gs *GeneralSettings) ProvisionInterval() int {
	return gs.provisionInterval
}

func (gs *GeneralSettings) ProvisionEnable() bool {
	return gs.provisionEnable
}

func (gs *GeneralSettings) HasChanged() bool {
	return gs.hasChanged
}

func (gs *GeneralSettings) UpdatedAt() time.Time {
	return gs.updatedAt
}

// Setters
func (gs *GeneralSettings) SetTimezone(timezone int) {
	gs.timezone = timezone
}
func (gs *GeneralSettings) SetTimeFormat(timeFormat string) {
	// На будующее нужно провалидировать формат
	switch timeFormat {
	case "12h", "24h":
		gs.timeFormat = timeFormat
	default:
		gs.timeFormat = "24h"
	}
}

func (gs *GeneralSettings) SetLanguage(language string) {
	switch language {
	case "en", "ru":
		gs.language = language
	default:
		gs.language = "ru"
	}
}

func (gs *GeneralSettings) SetWebPassword(webPassword string) {
	gs.webPassword = webPassword
}

func (gs *GeneralSettings) SetPhonePassword(phonePassword string) {
	gs.phonePassword = phonePassword
}

func (gs *GeneralSettings) EnableProvisionOnStart(provisionOnStart bool) {
	gs.provisionOnStart = provisionOnStart
}

func (gs *GeneralSettings) SetProvisionURL(provisionURL string) {
	gs.provisionURL = provisionURL
}

func (gs *GeneralSettings) SetProvisionInterval(provisionInterval int) {
	gs.provisionInterval = provisionInterval
}

func (gs *GeneralSettings) EnableProvisioning(provisionEnable bool) {
	gs.provisionEnable = provisionEnable
}

// Helpers
func (gs *GeneralSettings) markChanged(hasChanged bool) {
	gs.hasChanged = hasChanged
}

func (gs *GeneralSettings) MarkSynced() {
	gs.hasChanged = false
}

func (gs *GeneralSettings) setUpdatedAt(updatedAt time.Time) {
	gs.updatedAt = updatedAt
}
