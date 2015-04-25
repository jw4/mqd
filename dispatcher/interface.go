package dispatcher

type MailQueueCallbackFn func([]byte) bool

type MailQueueDispatcher interface {
	Process(MailQueueCallbackFn) error
}
