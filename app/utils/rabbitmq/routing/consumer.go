package routing

import (
	"ginskeleton/app/global/variable"
	"ginskeleton/app/utils/rabbitmq/error_record"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

func CreateConsumer(options ...OptionsConsumer) (*consumer, error) {

	conn, err := amqp.Dial(variable.ConfigYml.GetString("RabbitMq.Routing.Addr"))
	exchangeType := variable.ConfigYml.GetString("RabbitMq.Routing.ExchangeType")
	exchangeName := variable.ConfigYml.GetString("RabbitMq.Routing.ExchangeName")
	queueName := variable.ConfigYml.GetString("RabbitMq.Routing.QueueName")
	durable := variable.ConfigYml.GetBool("RabbitMq.Routing.Durable")
	reconnectInterval := variable.ConfigYml.GetDuration("RabbitMq.Routing.OffLineReconnectIntervalSec")
	retryTimes := variable.ConfigYml.GetInt("RabbitMq.Routing.RetryCount")

	if err != nil {
		return nil, err
	}

	cons := &consumer{
		connect:                     conn,
		exchangeType:                exchangeType,
		exchangeName:                exchangeName,
		queueName:                   queueName,
		durable:                     durable,
		connErr:                     conn.NotifyClose(make(chan *amqp.Error, 1)),
		offLineReconnectIntervalSec: reconnectInterval,
		retryTimes:                  retryTimes,
		receivedMsgBlocking:         make(chan struct{}),
		status:                      1,
	}

	for _, val := range options {
		val.apply(cons)
	}
	return cons, nil
}

type consumer struct {
	connect                     *amqp.Connection
	exchangeType                string
	exchangeName                string
	queueName                   string
	durable                     bool
	occurError                  error
	connErr                     chan *amqp.Error
	routeKey                    string
	callbackForReceived         func(receivedData string)
	offLineReconnectIntervalSec time.Duration
	retryTimes                  int
	callbackOffLine             func(err *amqp.Error)
	enableDelayMsgPlugin        bool
	receivedMsgBlocking         chan struct{}
	status                      byte
}

func (c *consumer) Received(routeKey string, callbackFunDealMsg func(receivedData string)) {
	defer func() {
		c.close()
	}()

	c.routeKey = routeKey
	c.callbackForReceived = callbackFunDealMsg

	go func(key string) {

		ch, err := c.connect.Channel()
		c.occurError = error_record.ErrorDeal(err)
		defer func() {
			_ = ch.Close()
		}()

		err = ch.ExchangeDeclare(
			c.exchangeName,
			c.exchangeType,
			c.durable,
			!c.durable,
			false,
			false,
			nil,
		)

		queue, err := ch.QueueDeclare(
			c.queueName,
			c.durable,
			true,
			false,
			false,
			nil,
		)
		c.occurError = error_record.ErrorDeal(err)

		err = ch.QueueBind(
			queue.Name,
			key,
			c.exchangeName,
			false,
			nil,
		)
		c.occurError = error_record.ErrorDeal(err)
		if err != nil {
			return
		}

		msgs, err := ch.Consume(
			queue.Name,
			"",
			true,
			false,
			false,
			false,
			nil,
		)
		c.occurError = error_record.ErrorDeal(err)
		if err == nil {
			for {
				select {
				case msg := <-msgs:

					if c.status == 1 && len(msg.Body) > 0 {
						callbackFunDealMsg(string(msg.Body))
					} else if c.status == 0 {
						return
					}
				}
			}
		} else {
			return
		}
	}(routeKey)

	if _, isOk := <-c.receivedMsgBlocking; isOk {
		c.status = 0
		close(c.receivedMsgBlocking)
	}

}

func (c *consumer) OnConnectionError(callbackOfflineErr func(err *amqp.Error)) {
	c.callbackOffLine = callbackOfflineErr
	go func() {
		select {
		case err := <-c.connErr:
			var i = 1
			for i = 1; i <= c.retryTimes; i++ {

				time.Sleep(c.offLineReconnectIntervalSec * time.Second)

				if c.status == 1 {
					c.receivedMsgBlocking <- struct{}{}
				}
				conn, err := CreateConsumer()
				if err != nil {
					continue
				} else {
					go func() {
						c.connErr = conn.connect.NotifyClose(make(chan *amqp.Error, 1))
						go conn.OnConnectionError(c.callbackOffLine)
						conn.Received(c.routeKey, c.callbackForReceived)
					}()

					if c.status == 0 {
						return
					}
					break
				}
			}
			if i > c.retryTimes {
				callbackOfflineErr(err)

				if c.status == 0 {
					return
				}
			}
		}
	}()
}

func (c *consumer) close() {
	_ = c.connect.Close()
}
