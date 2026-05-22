package topics

import (
	"ginskeleton/app/global/variable"
	amqp "github.com/rabbitmq/amqp091-go"
)

type OptionsProd interface {
	apply(*producer)
}

type OptionFunc func(*producer)

func (f OptionFunc) apply(prod *producer) {
	f(prod)
}

func SetProdMsgDelayParams(enableMsgDelayPlugin bool) OptionsProd {
	return OptionFunc(func(p *producer) {
		p.enableDelayMsgPlugin = enableMsgDelayPlugin
		p.exchangeType = "x-delayed-message"
		p.args = amqp.Table{
			"x-delayed-type": "topic",
		}
		p.exchangeName = variable.ConfigYml.GetString("RabbitMq.Topics.DelayedExchangeName")

		p.durable = true
	})
}

type OptionsConsumer interface {
	apply(*consumer)
}

type OptionsConsumerFunc func(*consumer)

func (f OptionsConsumerFunc) apply(cons *consumer) {
	f(cons)
}

func SetConsMsgDelayParams(enableDelayMsgPlugin bool) OptionsConsumer {
	return OptionsConsumerFunc(func(c *consumer) {
		c.enableDelayMsgPlugin = enableDelayMsgPlugin
		c.exchangeType = "x-delayed-message"
		c.exchangeName = variable.ConfigYml.GetString("RabbitMq.Topics.DelayedExchangeName")

		c.durable = true
	})
}
