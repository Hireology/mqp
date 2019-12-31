package main

import (
	"fmt"
	//"log"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	//"os"
	"testing"
)

// test BEHAVIOR not implementation
func setup() (conn *amqp.Connection) {
	connectionString := fmt.Sprintf("%s://%s:%s@%s:%d/%s",
		scheme, user, password, host, port, vhost)
	conn = connect(connectionString)
	return conn
}

func TestConnect(t *testing.T) {
	/*
		cases:
		- known good
		- known bad
	*/
	conn := setup()
	assert.Equal(t, conn.Config.Vhost, vhost)
	defer conn.Close()
}

func TestNewBasicPublishing(t *testing.T) {
	/*
	   cases:
	   - basic alphanumeric/special chars string
	   - unicode string
	   - emoji string
	*/
	pub := newBasicPublishing("hello world")
	assert.Equal(t, string(pub.Body), "hello world")
}

func TestPublishMessage(t *testing.T) {
	/*
	   cases:
	   - basic known good single message
	   - multiple messages
	*/
	conn := setup()
	channel, err := conn.Channel()
	assert.Nil(t, err)
	routingKey := "myRoutingKey"
	_, err = channel.QueueDeclare(
		routingKey, false, true, true, true, nil)
	assert.Nil(t, err)
	publishMessage(channel, "hello world")

	// got should be the last message in whatever test queue we used
	msg, _, err := channel.Get(routingKey, false)
	assert.Nil(t, err)
	assert.Equal(t, string(msg.Body), "hello world")
	defer conn.Close()
}

func TestProcessMessages(t *testing.T) {
	/*
	   cases:
	   - one message
	   - no messages
	   - 5 messages
	*/
	conn := setup()
	channel, err := conn.Channel()
	assert.Nil(t, err)

	routingKey := "myRoutingKey"
	_, err = channel.QueueDeclare(
		routingKey, false, true, true, true, nil)
	assert.Nil(t, err)

	publishMessage(channel, "ping")

	messages := messages(channel, routingKey, "myConsumer")
	processMessages(messages)
}

/*
func TestConnectTLS(t *testing.T) {
	//got := fmt.Sprintf("%+v", conn)
    got := "foo"
    want := "bar"
	assert.Equal(t, got, want)
}
*/
