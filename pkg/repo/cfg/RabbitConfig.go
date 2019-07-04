package cfg

type RabbitConfig interface {
	GetUrl() string
	GetEnvioDianQueue() QueueConfig
	GetNotificacionQueue() QueueConfig
}

type QueueConfig interface {
	GetName() string
}
