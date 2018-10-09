package nats

type NatsInterface interface {
	Close()
	PublishLogsCreated(mm LogsCreatedMessage) error
	SubscribeLogsCreated() (<-chan LogsCreatedMessage, error)
	OnLogsCreated(f func(LogsCreatedMessage)) error
}

var impl NatsInterface

func SetNatsInterface(es NatsInterface) {
	impl = es
}

func Close() {
	impl.Close()
}

func PublishLogsCreated(mm LogsCreatedMessage) error {
	return impl.PublishLogsCreated(mm)
}

func SubscribeLogsCreated() (<-chan LogsCreatedMessage, error) {
	return impl.SubscribeLogsCreated()
}

func OnLogsCreated(f func(LogsCreatedMessage)) error {
	return impl.OnLogsCreated(f)
}
