package app

type AdditionalLoggerModule struct{}

func (a *AdditionalLoggerModule) AdditionalLogger() bool {
	return true
}

type DefaultLoggerModule struct{}

func (a *DefaultLoggerModule) AdditionalLogger() bool {
	return false
}

type NoConfigModule struct{}

func (m *NoConfigModule) ConfigRequired() bool {
	return false
}

func (a *NoConfigModule) AdditionalLogger() bool {
	return false
}

type ConfigRequiredModule struct{}

func (m *ConfigRequiredModule) ConfigRequired() bool {
	return true
}
