package test

import (
	"fmt"
	"ginskeleton/app/global/my_errors"
	"ginskeleton/app/utils/rabbitmq/hello_world"
	"ginskeleton/app/utils/rabbitmq/publish_subscribe"
	"ginskeleton/app/utils/rabbitmq/routing"
	"ginskeleton/app/utils/rabbitmq/topics"
	"ginskeleton/app/utils/rabbitmq/work_queue"
	_ "ginskeleton/bootstrap"
	"os"
	"strconv"
	"testing"

	amqp "github.com/rabbitmq/amqp091-go"
)

func TestRabbitMqHelloWorldProducer(t *testing.T) {

	helloProducer, err := hello_world.CreateProducer()
	if err != nil {
		t.Errorf("HelloWorld单元测试未通过。%s\n", err.Error())
		os.Exit(1)
	}
	var res bool
	for i := 0; i < 10; i++ {
		str := fmt.Sprintf("%d_HelloWorld开始发送消息测试", i+1)
		res = helloProducer.Send(str)

	}

	helloProducer.Close()

	if res {
		t.Log("消息发送OK")
	} else {
		t.Errorf("HelloWorld模式消息发送失败")
	}
}

func TestMqHelloWorldConsumer(t *testing.T) {

	consumer, err := hello_world.CreateConsumer()
	if err != nil {
		t.Errorf("HelloWorld单元测试未通过。%s\n", err.Error())
		os.Exit(1)
	}

	consumer.OnConnectionError(func(err *amqp.Error) {
		t.Errorf(my_errors.ErrorsRabbitMqReconnectFail+"，%s\n", err.Error())
	})

	consumer.Received(func(receivedData string) {

		t.Logf("HelloWorld回调函数处理消息：--->%s\n", receivedData)
	})
}

func TestRabbitMqWorkQueueProducer(t *testing.T) {

	producer, _ := work_queue.CreateProducer()
	var res bool
	for i := 0; i < 10; i++ {
		str := fmt.Sprintf("%d_WorkQueue开始发送消息测试", i+1)
		res = producer.Send(str)

	}

	producer.Close()

	if res {
		t.Logf("消息发送OK")
	} else {
		t.Errorf("WorkQueue模式消息发送失败")
	}
}

func TestMqWorkQueueConsumer(t *testing.T) {

	consumer, err := work_queue.CreateConsumer()
	if err != nil {
		t.Errorf("WorkQueue单元测试未通过。%s\n", err.Error())
		os.Exit(1)
	}

	consumer.OnConnectionError(func(err *amqp.Error) {
		t.Errorf("%s, %s", my_errors.ErrorsRabbitMqReconnectFail, err.Error())
	})

	consumer.Received(func(receivedData string) {

		t.Logf("WorkQueue回调函数处理消息：--->%s\n", receivedData)
	})
}

func TestRabbitMqPublishSubscribeProducer(t *testing.T) {

	producer, err := publish_subscribe.CreateProducer()
	if err != nil {
		t.Errorf("WorkQueue 单元测试未通过。%s\n", err.Error())
		os.Exit(1)
	}
	var res bool
	for i := 0; i < 10; i++ {
		str := fmt.Sprintf("%d_PublishSubscribe开始发送消息测试", i+1)

		res = producer.Send(str, 1000)
		fmt.Println(str, res)

	}

	producer.Close()

	if res {
		t.Log("消息发送OK")
	} else {
		t.Errorf("PublishSubscribe 模式消息发送失败")
	}
}

func TestRabbitMqPublishSubscribeConsumer(t *testing.T) {

	consumer, err := publish_subscribe.CreateConsumer()
	if err != nil {
		t.Errorf("PublishSubscribe单元测试未通过。%s\n", err.Error())
		os.Exit(1)
	}

	consumer.OnConnectionError(func(err *amqp.Error) {
		t.Errorf("%s，%s\n", my_errors.ErrorsRabbitMqReconnectFail, err.Error())
	})

	consumer.Received(func(receivedData string) {

		t.Logf("PublishSubscribe回调函数处理消息：--->%s\n", receivedData)
	})
}

func TestRabbitMqRoutingProducer(t *testing.T) {

	producer, err := routing.CreateProducer()

	if err != nil {
		t.Errorf("Routing单元测试未通过。%s\n", err.Error())
		os.Exit(1)
	}
	var res bool
	var key string
	for i := 1; i <= 20; i++ {

		if i%2 == 0 {
			key = "key_even"
		} else {
			key = "key_odd"
		}

		res = producer.Send(key, strconv.Itoa(i)+"- Routing开始发送消息测试", 10000)

	}

	producer.Close()

	if res {
		t.Logf("消息发送OK")
	} else {
		t.Errorf("Routing 模式消息发送失败")
	}
}

func TestRabbitMqRoutingConsumer(t *testing.T) {
	consumer, err := routing.CreateConsumer()

	if err != nil {
		t.Errorf("Routing单元测试未通过。%s\n", err.Error())
		os.Exit(1)
	}

	consumer.OnConnectionError(func(err *amqp.Error) {
		t.Errorf("%s， %s\n", my_errors.ErrorsRabbitMqReconnectFail, err.Error())
	})

	consumer.Received("key_even", func(receivedData string) {
		fmt.Println("处理偶数的回调函数 ---> 收到消息内容: " + receivedData)

	})
}

func TestRabbitMqTopicsProducer(t *testing.T) {

	producer, err := topics.CreateProducer()
	if err != nil {
		t.Errorf("Routing单元测试未通过。%s\n", err.Error())
		os.Exit(1)
	}
	var res bool
	var key string
	for i := 1; i <= 10; i++ {

		if i%2 == 0 {
			key = "key.even"
		} else {
			key = "key.odd"
		}
		strData := fmt.Sprintf("%d_Routing_%s, 开始发送消息测试", i, key)

		res = producer.Send(key, strData, 10000)

	}

	producer.Close()

	if res {
		t.Logf("消息发送OK")
	} else {
		t.Errorf("topics 模式消息发送失败")
	}

}

func TestRabbitMqTopicsConsumer(t *testing.T) {
	consumer, err := topics.CreateConsumer()

	if err != nil {
		t.Errorf("Routing单元测试未通过。%s\n", err.Error())
		os.Exit(1)
	}

	consumer.OnConnectionError(func(err *amqp.Error) {
		t.Errorf("%s， %s\n", my_errors.ErrorsRabbitMqReconnectFail, err.Error())
	})

	consumer.Received("#.odd", func(receivedData string) {

		t.Logf("模糊匹配偶数键：--->%s\n", receivedData)
	})
}
