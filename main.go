package main

import (
	// packages
	//"fmt"
	"github.com/streadway/amqp"
	"log"
	"time"
	//"github.com/Hireology/go_shared/messaging"
)

var (
	scheme   = "amqp"
	user     = "mqp"
	password = "mqptest"
	host     = "127.0.0.1"
	port     = 5672
	//vhost    = "/"
	vhost = "mqp"
)

func main() {
	// stuff
}

func failOnError(prefix string, err error) {
	if err != nil {
		log.Fatalf("%s: %s", prefix, err)
	}
}

func connect(connectionString string) (conn *amqp.Connection) {
	conn, err := amqp.Dial(connectionString)
	failOnError("connection failure", err)

	return conn
}

/*
// connect with tls
func connectTLS() {
}
*/

func newBasicPublishing(message string) *amqp.Publishing {
	publishing := amqp.Publishing{
		DeliveryMode: amqp.Transient,
		Timestamp:    time.Now(),
		ContentType:  "text/plain",
		Body:         []byte(message),
	}
	return &publishing
}

// publish a message
func publishMessage(c *amqp.Channel, message string) amqp.Publishing {
	exchange := ""
	routingKey := "myRoutingKey"
	publishing := newBasicPublishing(message)

	err := c.Publish(exchange, routingKey, true, false, *publishing)
	// TODO where does this prefix come from?
	failOnError("basic.publish", err)
	//defer c.Close()
	return *publishing
}

func messages(c *amqp.Channel, queue, consumer string) <-chan amqp.Delivery {
	deliver, err := c.Consume(
		queue,
		consumer,
		true,  // autoAck
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,   // args
	)
	//c.Close()
	failOnError("channel.consume", err)
	return deliver
}

func processMessages(messages <-chan amqp.Delivery) {
	for d := range messages {
		log.Println("received:", string(d.Body))
	}
}
