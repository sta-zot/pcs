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
// =========================================================
func (ss *SipSettings) Extension() int      { return ss.extension }
func (ss *SipSettings) Secret() string      { return ss.secret }
func (ss *SipSettings) DisplayName() string { return ss.displayName }
func (ss *SipSettings) ServerHost() string  { return ss.serverHost }
func (ss *SipSettings) ServerPort() int     { return ss.serverPort }
func (ss *SipSettings) Transport() string   { return ss.transport }
func (ss *SipSettings) DtmfMode() string    { return ss.dtmfMode }
func (ss *SipSettings) ExpireTime() int     { return ss.expireTime }

// Setters
// =========================================================
func (ss *SipSettings) SetDisplayName(displayName string) {
	ss.displayName = displayName
	ss.MarkChanged()
}
func (ss *SipSettings) SetServerHost(serverHost string) {
	ip, err := validateIP(serverHost)
	if err == nil {
		ss.serverHost = ip
		ss.MarkChanged()
		return
	}

	if validateDN(serverHost) {
		ss.serverHost = serverHost
		ss.MarkChanged()
		return
	}

}
func (ss *SipSettings) SetServerPort(serverPort int) {
	if serverPort > 0 && serverPort <= 65535 {
		ss.serverPort = serverPort
		ss.touch()
		return
	}
}
func (ss *SipSettings) SetTransport(transport string) {
	ss.transport = transport
	ss.MarkChanged()
}
func (ss *SipSettings) SetDtmfMode(dtmfMode string) {
	ss.dtmfMode = dtmfMode
	ss.MarkChanged()
}
func (ss *SipSettings) SetExpireTime(expireTime int) {
	ss.expireTime = expireTime
	ss.MarkChanged()
}

// HasChanged возвращает true, если настройки были изменены
func (ss *SipSettings) HasChanged() bool {
	return ss.hasChanged
}

func (ss *SipSettings) MarkChanged() {
	ss.hasChanged = true
	ss.touch()
}
func (ss *SipSettings) touch() {
	ss.updatedAt = time.Now()
}
