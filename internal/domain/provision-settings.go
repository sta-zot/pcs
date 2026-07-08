package domain

import "time"

type ProvisioningSettings struct {
	enabled     bool      // Включено/выключено
	url         string    // URL для конфигурации
	username    string    // Имя пользователя
	password    string    // Пароль
	schedule    string    // Расписание проверки обновлений
	firmwareURL string    // URL для прошивки
	lastSeen    time.Time // время последнего обновления
	hasChanged  bool      // для точечного обновления
	updatedAt   time.Time // время последнего обновления

}
