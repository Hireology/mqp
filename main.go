package main

import (
	// packages
	"flag"
	//"fmt"
	"github.com/streadway/amqp"
	"log"
	"time"
	//"github.com/Hireology/go_shared/messaging"
)

/*
var (
	scheme   = "amqp"
	user     = "mqp"
	password = "mqptest"
	host     = "127.0.0.1"
	port     = 5672
	vhost    = "mqp" // leading slash is not needed
)
*/

const (
	DefaultScheme   = "amqp"
	DefaultUser     = "mqp"
	DefaultPassword = "mqptest"
	DefaultHost     = "127.0.0.1"
	DefaultPort     = 5672
	DefaultVhost    = "mqp"
)

// TODO when to create one of these?
// TODO do we actually connect when the object is created or just have an empty ref?
type MQ struct {
	URI        *amqp.URI
	Connection *amqp.Connection
}

// NewMQ initializes RabbitMQ and returns a pointer to it
func NewMQ(uri *amqp.URI) (*MQ, error) {
	mq := MQ{
		URI: uri,
	}
	err := mq.connect()
	return &mq, err
}

func main() {
	// stuff
}

func parseFlags() *amqp.URI {
	// TODO flags won't always just be URI parameters
	var uri = amqp.URI{
		Scheme:   *flag.String("scheme", DefaultScheme, "Connection scheme"),
		Host:     *flag.String("host", DefaultHost, "Connection host"),
		Port:     *flag.Int("port", DefaultPort, "Connection port"),
		Username: *flag.String("user", DefaultUser, "Connection user"),
		Password: *flag.String("password", DefaultPassword, "Connection password"),
		Vhost:    *flag.String("vhost", DefaultVhost, "Connection vhost"),
	}
	flag.Parse()
	return &uri
}

func failOnError(prefix string, err error) {
	if err != nil {
		log.Fatalf("%s: %s", prefix, err)
	}
}

func (mq *MQ) connect() error {
	conn, err := amqp.Dial(mq.URI.String())
	//failOnError("connection failure", err)
	mq.Connection = conn

	return err
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
	failOnError("channel.consume", err)
	return deliver
}

func processMessages(messages <-chan amqp.Delivery) {
	for d := range messages {
		log.Println("received:", string(d.Body))
	}
}
