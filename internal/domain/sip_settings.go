package domain

import "time"

type SipSettings struct {
	extension   int       // extension name / ID (например, "101")
	secret      string    // Пароль из sip/pjsip settings
	displayName string    // Отображаемое имя (CallerID Name)
	serverHost  string    // IP/Хост АТС
	serverPort  int       // 5060 / 5061
	transport   string    // UDP, TCP, TLS
	dtmfMode    string    // rfc2833, info, inband
	expireTime  int       // Тайм-аут регистрации
	updatedAt   time.Time // время последнего обновления
	hasChanged  bool      // для точечного обновления
}

func NewSipSettings(
	extension int,
	secret string,
	serverHost string,

) *SipSettings {
	return &SipSettings{
		extension:  extension,
		secret:     secret,
		serverHost: serverHost,
	}
}

// Getters
// Если поребуется логика приобразования можно будет включить без проблем
//=========================================================

// Setters
// =========================================================

func (ss *SipSettings) SetDisplayName(displayName string) { ss.displayName = displayName }
func (ss *SipSettings) SetServerHost(serverHost string)   {
    ip, err := validateIP(serverHost)
    if err == nil{
        return
    }
    if
    ss.serverHost = serverHost
}
