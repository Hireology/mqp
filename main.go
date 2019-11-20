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

	q, err := c.QueueDeclare(routingKey, false, true, true, true, nil)
	failOnError("queue declare", err)
	err = c.Publish(exchange, q.Name, true, false, *publishing)
	// TODO where does this prefix come from?
	failOnError("basic.publish", err)
	defer c.Close()
	return *publishing
}

// consume a message
func consumeMessage(c *amqp.Channel, queue, consumer string) amqp.Delivery {
	//func (ch *Channel) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args Table) (<-chan Delivery, error)
	deliver, err := c.Consume(queue, consumer, false, false, false, false, nil)
	failOnError("channel.consume", err)
	// what should this actually return?
	// what do we get from <-chan amqp.Deliver ?
	for {
		msg, ok := <-deliver
		if ok == false {
			log.Println("channel closed")
			break
		}
		log.Println("received:", msg, ok)
		err := msg.Ack(true)
		failOnError("delivery.ack", err)
	}

	stuff := amqp.Delivery{}
	return stuff
}
