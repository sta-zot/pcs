1. Агрегат в домене отвечающий за все параметры телефона

- PhoneConfig
```go
type PhoneConfig struct {
	id int
	macAddr string
	vendor string
	model string
	userAgent string
	isComplete bool
	hasChanged bool
	lastSeen time.Time
	changedAt time.Time
	updatedAt time.Time

	// дальше агрегация 
	network NetworkSettings
	sip SipSettings
	provisioning ProvisioningSettings
	general GeneralSettings
	lines []LineSettings
}
```

-
