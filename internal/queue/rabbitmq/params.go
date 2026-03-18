package rabbitmq

const (
	ExchangeName       = "notifications"
	ExchangeType       = "direct"
	ExchangeKind       = "x-delayed-message"
	ExchangeDurable    = true
	ExchangeAutoDelete = false
	ExchangeInternal   = false
)

const (
	DeadLetterExchangeName       = "dead-notifications"
	DeadLetterExchangeKind       = "fanout"
	DeadLetterExchangeDurable    = true
	DeadLetterExchangeAutoDelete = false
	DeadLetterExchangeInternal   = false
)

const (
	EmailQueueName  = "email queue"
	EmailRoutingKey = "email"
	TgQueueName     = "telegram queue"
	TgRoutingKey    = "telegram"
	QueueDurable    = true
	QueueAutoDelete = false

	BingNoWait  = false
	ContentType = "application/json"
)

const (
	DeadLetterQueueName       = "dead letter queue"
	DeadLetterRoutingKey      = ""
	DeadLetterQueueDurable    = true
	DeadLetterQueueAutoDelete = false
)
