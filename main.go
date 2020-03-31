package main

import (
	"flag"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"os"
	"time"
)

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

// newMQ initializes RabbitMQ and returns a pointer to it
func newMQ(uri *amqp.URI) (*MQ, error) {
	log.Printf("connecting to %s:%d", uri.Host, uri.Port)
	mq := MQ{
		URI: uri,
	}
	err := mq.connect()
	return &mq, err
}

func (mq *MQ) setupChannel(routingKey string) (*amqp.Channel, error) {
	conn := mq.Connection
	channel, err := conn.Channel()

	_, err = channel.QueueDeclare(
		routingKey, false, true, true, true, nil)

	return channel, err
}

func publishMessages(channel *amqp.Channel, msg string, expectedMessages int) {
	for i := 1; i <= expectedMessages; i++ {
		msg := fmt.Sprintf("%s: %d of %d", msg, i, expectedMessages)
		publishMessage(channel, msg)
		log.Printf("published: %s", msg)
	}
}

func channelTimeout(channel Closer, timeoutSeconds int) {
	time.Sleep(time.Duration(timeoutSeconds) * time.Second)
	channel.Close()
}

type Closer interface {
	Close() error
}

func main() {
	// connection setup
	expectedMessages := 5
	routingKey := "myRoutingKey"
	consumer := "myConsumer"
	uri := parseFlags()

	mq, err := newMQ(uri)
	// TODO handle error
	failOnError("connection", err)

	channel, err := mq.setupChannel(routingKey)
	// TODO handle error
	failOnError("channel setup", err)

	publishMessages(channel, "hello world", expectedMessages)

	go channelTimeout(channel, 1)
	messages := messages(channel, routingKey, consumer)
	processMessages(messages)

	// TODO allow user to specify N
	// TODO display statistics
	/* exit codes
	0 = all messages published and delivered successfully
	1 = at least one message NACK'd
	2 = any kind of failure (messages dropped, unaccounted for, connection failed, etc)
	*/
	mq.Connection.Close()

	// exit with relevant error code
	os.Exit(0)
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
	// TODO pass error up
	failOnError("basic.publish", err)
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

// tell processMessages about how many messages it should get?
func processMessages(messages <-chan amqp.Delivery) {
	for d := range messages {
		log.Println("received:", string(d.Body))
	}
}
