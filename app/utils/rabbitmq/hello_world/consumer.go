package hello_world

import (
	"ginskeleton/app/global/variable"
	"ginskeleton/app/utils/rabbitmq/error_record"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

func CreateConsumer() (*consumer, error) {

	conn, err := amqp.Dial(variable.ConfigYml.GetString("RabbitMq.HelloWorld.Addr"))
	queueName := variable.ConfigYml.GetString("RabbitMq.HelloWorld.QueueName")
	durable := variable.ConfigYml.GetBool("RabbitMq.HelloWorld.Durable")
	chanNumber := variable.ConfigYml.GetInt("RabbitMq.HelloWorld.ConsumerChanNumber")
	reconnectInterval := variable.ConfigYml.GetDuration("RabbitMq.HelloWorld.OffLineReconnectIntervalSec")
	retryTimes := variable.ConfigYml.GetInt("RabbitMq.HelloWorld.RetryCount")

	if err != nil {

		return nil, err
	}
	cons := &consumer{
		connect:                     conn,
		queueName:                   queueName,
		durable:                     durable,
		chanNumber:                  chanNumber,
		connErr:                     conn.NotifyClose(make(chan *amqp.Error, 1)),
		offLineReconnectIntervalSec: reconnectInterval,
		retryTimes:                  retryTimes,
		receivedMsgBlocking:         make(chan struct{}),
		status:                      1,
	}
	return cons, nil
}

type consumer struct {
	connect                     *amqp.Connection
	queueName                   string
	durable                     bool
	chanNumber                  int
	occurError                  error
	connErr                     chan *amqp.Error
	callbackForReceived         func(receivedData string)
	offLineReconnectIntervalSec time.Duration
	retryTimes                  int
	callbackOffLine             func(err *amqp.Error)
	receivedMsgBlocking         chan struct{}
	status                      byte
}

func (c *consumer) Received(callbackFunDealSmg func(receivedData string)) {
	defer func() {
		c.close()
	}()

	c.callbackForReceived = callbackFunDealSmg

	for i := 1; i <= c.chanNumber; i++ {
		go func(chanNo int) {
			ch, err := c.connect.Channel()
			c.occurError = error_record.ErrorDeal(err)
			defer func() {
				_ = ch.Close()
			}()

			queue, err := ch.QueueDeclare(
				c.queueName,
				c.durable,
				true,
				false,
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
							callbackFunDealSmg(string(msg.Body))
						} else if c.status == 0 {
							return
						}
					}
				}
			} else {
				return
			}
		}(i)
	}

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
						conn.Received(c.callbackForReceived)
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
