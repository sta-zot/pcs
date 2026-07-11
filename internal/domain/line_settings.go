package domain

import (
	"time"
)

// Структура описывет настройки линии на SIP телефоне
type LineSettings struct {
	id          int  // ID линии
	slot        int  // Номер линии
	enabled     bool // Включена/выключена хотя на мой взгляд лишнее
	extensionID int  // Идентификатор SIP аккаунта в БД asterisk для быстрого поиска и запроса требуемых настроек
	// extension Label и DisplayName и пароль берётся из БД asterisk таблицы sip(поля extension, name, )
	SipInfo    SipSettings
	hasChanged bool      // для точечного обновления
	createdAt  time.Time // время создания
	updatedAt  time.Time // время последнего обновления
}

// ===========================================
// СОЗДАНИЕ СТРУКТУРЫ
// ===========================================
func NewLineSettings(slot, extId int) *LineSettings {
	if slot < 1 {
		slot = 1
	}
	if slot > 10 {
		slot = 10
	}
	return &LineSettings{
		extensionID: extId,
		slot:        slot,
	}
}

// ============================================
// ГЕТТЕРЫ
// ============================================

func (l *LineSettings) ID() int              { return l.id }
func (l *LineSettings) Slot() int            { return l.slot }
func (l *LineSettings) Enabled() bool        { return l.enabled }
func (l *LineSettings) ExtensionID() int     { return l.extensionID }
func (l *LineSettings) Sip() *SipSettings    { return &l.SipInfo }
func (l *LineSettings) HasChanged() bool     { return l.hasChanged }
func (l *LineSettings) CreatedAt() time.Time { return l.createdAt }
func (l *LineSettings) UpdatedAt() time.Time { return l.updatedAt }

// ============================================
// СЕТТЕРЫ
// ============================================

func (l *LineSettings) SetID(id int) {
	l.id = id
}
func (l *LineSettings) SetSipInfo(
	extension int,
	secret string,
	displayName string,
	serverHost string,
) {
	// используем сетеры SipInfo
}
func (l *LineSettings) SetSlot(slot int) error {
	if slot < 0 {
		return ErrInvalidSlot
	}
	l.slot = slot
	l.MarkChanged()
	return nil
}
func (l *LineSettings) SetEnabled(enabled bool) {
	l.enabled = enabled
	l.touch()
	l.MarkChanged()
}

func (l *LineSettings) SetCreatedAt(t time.Time) {
	l.createdAt = t
} // Нужно при заполнении из БД

func (l *LineSettings) SetUpdatedAt(t time.Time) {
	l.updatedAt = t
}

// ============================================
// БИЗНЕС-МЕТОДЫ
// ============================================

func (l *LineSettings) MarkChanged() {
	l.hasChanged = true

}

func (l *LineSettings) MarkSynced() {
	l.hasChanged = false
	l.touch()
}

func (l *LineSettings) touch() {
	l.updatedAt = time.Now()

}
