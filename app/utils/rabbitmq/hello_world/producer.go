package hello_world

import (
	"ginskeleton/app/global/variable"
	"ginskeleton/app/utils/rabbitmq/error_record"
	amqp "github.com/rabbitmq/amqp091-go"
)

func CreateProducer() (*producer, error) {

	conn, err := amqp.Dial(variable.ConfigYml.GetString("RabbitMq.HelloWorld.Addr"))
	queueName := variable.ConfigYml.GetString("RabbitMq.HelloWorld.QueueName")
	dura := variable.ConfigYml.GetBool("RabbitMq.HelloWorld.Durable")

	if err != nil {
		variable.ZapLog.Error(err.Error())
		return nil, err
	}

	prod := &producer{
		connect:   conn,
		queueName: queueName,
		durable:   dura,
	}
	return prod, nil
}

type producer struct {
	connect    *amqp.Connection
	queueName  string
	durable    bool
	occurError error
}

func (p *producer) Send(data string) bool {

	ch, err := p.connect.Channel()
	p.occurError = error_record.ErrorDeal(err)

	defer func() {
		_ = ch.Close()
	}()

	_, err = ch.QueueDeclare(
		p.queueName,
		p.durable,
		!p.durable,
		false,
		false,
		nil,
	)
	p.occurError = error_record.ErrorDeal(err)

	msgPersistent := amqp.Transient
	if p.durable {
		msgPersistent = amqp.Persistent
	}

	if err == nil {
		err = ch.Publish(
			"",
			p.queueName,
			false,
			false,
			amqp.Publishing{
				DeliveryMode: msgPersistent,
				ContentType:  "text/plain",
				Body:         []byte(data),
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
