package data

type Consumer interface {
	ConsumeValue(item Item, value string)
	ConsumeError(item Item, err error)
}
