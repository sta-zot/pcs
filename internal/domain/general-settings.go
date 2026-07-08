package domain

import "time"

type GeneralSettings struct {
	timezone      string    // Часовой пояс
	timeFormat    string    // Формат времени (12h/24h)
	language      string    // Язык интерфейса
	ringVolume    int       // Громкость звонка
	speakerVolume int       // Громкость динамика
	webPassword   string    // Пароль для веб-интерфейса
	phonePassword string    // Пароль для телефона
	hasChanged    bool      // для точечного обновления
	updatedAt     time.Time // время последнего обновления
}
