package domain

import (
	"fmt"
	"time"
)

// Структура описывет настройки линии на SIP телефоне
type LineSettings struct {
	id          int  // ID линии
	slot        int  // Номер линии
	enabled     bool // Включена/выключена хотя на мой взгляд лишнее
	extensionID int  // Идентификатор SIP аккаунта в БД asterisk для быстрого поиска и запроса требуемых настроек
	// extension Label и DisplayName и пароль берётся из БД asterisk таблицы sip(поля extension, name, )
	isComplete bool // Определяет, завершены ли все настройки линии
	SipInfo    SipSettings
	source     string    // источник данных
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
func (l *LineSettings) Source() string       { return l.source }

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

func (l *LineSettings) SetSource(source string) {
	l.source = source
}

// ============================================
//
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

func (l *LineSettings) Complete() (bool, error) {
	err := l.checkComplete()
	return l.isComplete, err
}

func (l *LineSettings) checkComplete() error {
	if l.SipInfo.secret == "" {
		l.isComplete = false
		return fmt.Errorf("extension secret is empty")
	}
	if l.SipInfo.serverHost == "" {
		l.isComplete = false
		return fmt.Errorf("server host is empty")
	}
	l.isComplete = true
	return nil
}
