package rabbitmq

type QueueInfo struct {
	MessagesCount int `json:"messages_persistent"`
}
