package app

type AdditionalLoggerModule struct{}

func (a *AdditionalLoggerModule) AdditionalLogger() bool {
	return true
}

type DefaultLoggerModule struct{}

func (a *DefaultLoggerModule) AdditionalLogger() bool {
	return false
}
