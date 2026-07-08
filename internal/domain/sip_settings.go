package domain

import "time"

type SipSettings struct {
	server        string    // SIP сервер
	port          int       // Порт
	transport     int       // Транспорт (UDP/TCP/TLS)
	outBoundProxy string    // Outbound proxy
	outboundPort  int       // Outbound port
	stunServer    string    // STUN сервер
	updatedAt     time.Time // время последнего обновления
	hasChanged    bool      // для точечного обновления
}
