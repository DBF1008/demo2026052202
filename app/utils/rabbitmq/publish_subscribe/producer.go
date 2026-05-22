package publish_subscribe

import (
	"ginskeleton/app/global/variable"
	"ginskeleton/app/utils/rabbitmq/error_record"
	amqp "github.com/rabbitmq/amqp091-go"
)

func CreateProducer(options ...OptionsProd) (*producer, error) {

	conn, err := amqp.Dial(variable.ConfigYml.GetString("RabbitMq.PublishSubscribe.Addr"))
	exchangeType := variable.ConfigYml.GetString("RabbitMq.PublishSubscribe.ExchangeType")
	exchangeName := variable.ConfigYml.GetString("RabbitMq.PublishSubscribe.ExchangeName")
	queueName := variable.ConfigYml.GetString("RabbitMq.PublishSubscribe.QueueName")
	durable := variable.ConfigYml.GetBool("RabbitMq.PublishSubscribe.Durable")

	if err != nil {
		variable.ZapLog.Error(err.Error())
		return nil, err
	}

	prod := &producer{
		connect:      conn,
		exchangeType: exchangeType,
		exchangeName: exchangeName,
		queueName:    queueName,
		durable:      durable,
		args:         nil,
	}

	for _, val := range options {
		val.apply(prod)
	}
	return prod, nil
}

type producer struct {
	connect              *amqp.Connection
	exchangeType         string
	exchangeName         string
	queueName            string
	durable              bool
	occurError           error
	enableDelayMsgPlugin bool
	args                 amqp.Table
}

func (p *producer) Send(data string, delayMillisecond int) bool {

	ch, err := p.connect.Channel()
	p.occurError = error_record.ErrorDeal(err)
	defer func() {
		_ = ch.Close()
	}()

	err = ch.ExchangeDeclare(
		p.exchangeName,
		p.exchangeType,
		p.durable,
		!p.durable,
		false,
		false,
		p.args,
	)
	p.occurError = error_record.ErrorDeal(err)

	msgPersistent := amqp.Transient
	if p.durable {
		msgPersistent = amqp.Persistent
	}

	if err == nil {
		err = ch.Publish(
			p.exchangeName,
			p.queueName,
			false,
			false,
			amqp.Publishing{
				DeliveryMode: msgPersistent,
				ContentType:  "text/plain",
				Body:         []byte(data),
				Headers: amqp.Table{
					"x-delay": delayMillisecond,
				},
			})
	}

	p.occurError = error_record.ErrorDeal(err)
	if p.occurError != nil {
		return false
	} else {
		return true
	}
}

func (p *producer) Close() {
	_ = p.connect.Close()
}
