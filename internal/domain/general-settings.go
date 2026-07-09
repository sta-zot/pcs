package domain

import "time"

type GeneralSettings struct {
	timezone          string    // Часовой пояс
	timeFormat        string    // Формат времени (12h/24h)
	language          string    // Язык интерфейса
	ringVolume        int       // Громкость звонка
	speakerVolume     int       // Громкость динамика
	webPassword       string    // Пароль для веб-интерфейса
	phonePassword     string    // Пароль для телефона
	provisionOnStart  bool      // static.auto_provision.power_on
	provisionURL      string    // static.auto_provision.server.url =
	provisionInterval int       // static.auto_provision.flexible.interval =
	provisionEnable   bool      //static.auto_provision.flexible.enable =
	hasChanged        bool      // для точечного обновления
	updatedAt         time.Time // время последнего обновления
}
