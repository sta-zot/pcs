package domain

import (
	"time"
)

// Структура описывет настройки линии на SIP телефоне
type LineSettings struct {
	id          int  // ID линии
	lineNumber  int  // Номер линии
	enabled     bool // Включено/выключено
	extensionID int  // Идентификатор SIP аккаунта в БД asterisk
	// extension Label и DisplayName и пароль берётся из БД asterisk таблицы Users(поля extension, name, )
	codecs     []string  // Кодеки: ["PCMU","PCMA","G722"] из строки бд
	srtp       bool      // SRTP: true/false
	hasChanged bool      // для точечного обновления
	createdAt  time.Time // время создания
	updatedAt  time.Time // время последнего обновления
}

// ===========================================
// СОЗДАНИЕ СТРУКТУРЫ
// ===========================================
func NewLineSettings(lineNum int) *LineSettings {
	if lineNum < 1 {
		lineNum = 1
	}
	if lineNum > 10 {
		lineNum = 10
	}
	return &LineSettings{
		lineNumber: lineNum,
	}
}

// ============================================
// ГЕТТЕРЫ
// ============================================

func (l *LineSettings) ID() int              { return l.id }
func (l *LineSettings) LineNumber() int      { return l.lineNumber }
func (l *LineSettings) Enabled() bool        { return l.enabled }
func (l *LineSettings) ExtensionID() int     { return l.extensionID }
func (l *LineSettings) Codecs() []string     { return l.codecs }
func (l *LineSettings) SRTP() bool           { return l.srtp }
func (l *LineSettings) HasChanged() bool     { return l.hasChanged }
func (l *LineSettings) CreatedAt() time.Time { return l.createdAt }
func (l *LineSettings) UpdatedAt() time.Time { return l.updatedAt }

// ============================================
// СЕТТЕРЫ
// ============================================

func (l *LineSettings) SetID(id int) {
	l.id = id
}

func (l *LineSettings) SetLineNumber(num int) {
	l.lineNumber = num
	l.touch()
}

func (l *LineSettings) SetEnabled(enabled bool) {
	l.enabled = enabled
	l.touch()
}

func (l *LineSettings) SetExtensionID(id int) {
	l.extensionID = id
	l.touch()
}

func (l *LineSettings) SetCodecs(codecs []string) {
	l.codecs = codecs
	l.touch()
}

func (l *LineSettings) SetSRTP(srtp bool) {
	l.srtp = srtp
	l.touch()
}

func (l *LineSettings) SetHasChanged(changed bool) {
	l.hasChanged = changed
	// Не вызываем touch() чтобы избежать рекурсии
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
	if !l.hasChanged {
		l.MarkChanged()
	}

}
