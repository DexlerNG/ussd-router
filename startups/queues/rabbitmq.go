package queues

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"os"
	"strconv"
	"sync"
)

const CHARGE_AUTHORIZATION_PROCESSING_QUEUE = "telco-subscribe-billing.authorization.processing"
const SUBSCRIPTION_PROCESSING_QUEUE = "telco-subscribe-billing.subscribe.processing"
const DATA_SYNC_PROCESSING_QUEUE = "telco-subscribe-billing.data-sync.processing"
const CHARGE_PROCESSING_QUEUE = "telco-subscribe-billing.charge.processing"
const DEFAULT_EXCHANGE = "subscription-billing-engine.exchange"
const DEAD_LETTER_QUEUE = "telco-subscribe-billing.incoming.dead.letter"



var channel *amqp.Channel
var connection *amqp.Connection
var rabbitConnectOnce sync.Once

type Fn func([]byte) error

func GetRabbitMQClient() *amqp.Channel {

	rabbitConnectOnce.Do(func() {
		println("Calling RabbitMQ .Once ......", os.Getenv("RABBITMQ_URL") != "")
		// Connect to the rabbitMQ instance
		createdConnection, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))

		if err != nil {
			panic("could not establish connection with RabbitMQ:" + err.Error())
		}
		fmt.Println("Creating RabbitMQ Channel")

		// Create a channel from the connection. We'll use channels to access the data in the queue rather than the connection itself.
		createdChannel, err := createdConnection.Channel()
		if err != nil {
			panic("could not open channel with RabbitMQ:" + err.Error())
		}
		channel = createdChannel
		connection = createdConnection
		//init Queues and Exchanges
		//err = channel.ExchangeDeclare(".exchange", "fanout", true, false, false, false, nil)
		//if err != nil {
		//	fmt.Println("Error Creating exchange", err)
		//}

		queue, err := channel.QueueDeclare(CHARGE_AUTHORIZATION_PROCESSING_QUEUE, true, false, false, false, nil)
		fmt.Println("Queue", queue)
		if err != nil {
			fmt.Println("Error Creating Queue: CHARGE_AUTHORIZATION_PROCESSING_QUEUE", err)
		}

		queue, err = channel.QueueDeclare(CHARGE_PROCESSING_QUEUE, true, false, false, false, nil)
		fmt.Println("Queue", queue)
		if err != nil {
			fmt.Println("Error Creating Queue: CHARGE_PROCESSING_QUEUE", err)
		}

		queue, err = channel.QueueDeclare(SUBSCRIPTION_PROCESSING_QUEUE, true, false, false, false, nil)
		fmt.Println("Queue", queue)
		if err != nil {
			fmt.Println("Error Creating Queue: SUBSCRIPTION_PROCESSING_QUEUE", err)
		}


		queue, err = channel.QueueDeclare(DATA_SYNC_PROCESSING_QUEUE, true, false, false, false, nil)
		fmt.Println("Queue", queue)
		if err != nil {
			fmt.Println("Error Creating Queue: DATA_SYNC_PROCESSING_QUEUE", err)
		}


		err = channel.ExchangeDeclare(DEFAULT_EXCHANGE, "fanout", true, false, false, false, nil)
		if err != nil {
			fmt.Println("Error Creating Exchange: DEFAULT_EXCHANGE", err)
		}
	})
	//fmt.Println("channel", channel)
	return channel
}

func GetRabbitMQConnection() *amqp.Connection {
	return connection
}

func RabbitMQPublishToExchange(exchangeName string, data interface{}) error {

	jsonPayload, _ := json.Marshal(data)
	//fmt.Println("QueuePayload", jsonPayload)
	message := amqp.Publishing{
		Body: []byte(string(jsonPayload)),
	}
	err := GetRabbitMQClient().Publish(exchangeName, "", false, false, message)
	if err != nil {
		fmt.Println("Error While Pushing to Exchange: ", err)
		return err
	}
	return nil
}

func RabbitMQPublishToQueue(queueName string, data interface{}) error {

	jsonPayload, _ := json.Marshal(data)
	fmt.Println("QueuePayload", jsonPayload, queueName)
	message := amqp.Publishing{
		Body: []byte(string(jsonPayload)),
	}
	err := GetRabbitMQClient().Publish("", queueName, false, false, message)
	if err != nil {
		fmt.Println("Error While Pushing to Queue: ", err)
		return err
	}

	return nil
}

func Listen(listenMap map[string]Fn) {
	prefetch, err := strconv.Atoi(os.Getenv("QUEUE_PREFETCH"))
	if err != nil{
		prefetch = 1
	}
	for queueName, functionToCall := range listenMap {
		go individualListen(queueName, functionToCall, prefetch)
	}

}

func individualListen(queueName string, functionToCall Fn, prefetch int){

	channel.Qos(
		prefetch,           // prefetch count
		0,               // prefetch size
		false,           // global
	)
	msgs, err := GetRabbitMQClient().Consume(
		queueName,     // queue
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		panic("Failed to register a consumer")
	}
	forever := make(chan bool)
	func() {
		for d := range msgs {
			fmt.Printf("Received a message: %s\n", string(d.Body))
			err = functionToCall(d.Body)
			if err != nil {
				fmt.Printf("Error Processing Message: %s\n", d.Body)
				//push to dead letter Queue,
				//Map result to slice
				jsonPayload, _ := json.Marshal(&map[string]interface{}{
					"reason":  err.Error(),
					"error":   err.Error(),
					"queue":   queueName,
					"payload": d.Body,
				})
				//fmt.Println("QueuePayload", jsonPayload)
				message := amqp.Publishing{
					Body: []byte(string(jsonPayload)),
				}
				_ = GetRabbitMQClient().Publish("", DEAD_LETTER_QUEUE, false, false, message)
				_ = d.Ack(true)
				continue
			}
			_ = d.Ack(true)
			fmt.Println("Message Acknowledge")
		}
	}()

	fmt.Println(" [*] Waiting for messages 0n "+ queueName + ". To exit press CTRL+C")
	<-forever
}