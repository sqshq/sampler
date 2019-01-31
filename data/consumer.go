package data

type Consumer interface {
	ConsumeValue(value string, label string)
	ConsumeError(err error)
}
