package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type MQPTestSuite struct {
	suite.Suite
	mq *MQ
}

func (suite *MQPTestSuite) SetupSuite() {
	uri := parseFlags()
	suite.mq, _ = NewMQ(uri)
}

func TestMQPTestSuite(t *testing.T) {
	suite.Run(t, new(MQPTestSuite))
}

func (suite *MQPTestSuite) TestConnect() {
	/*
		cases:
		- known good
		- known bad
	*/
	conn := suite.mq.Connection
	assert.Equal(suite.T(), conn.Config.Vhost, "mqp")
	//defer conn.Close()
}

func (suite *MQPTestSuite) TestNewBasicPublishing() {
	/*
	   cases:
	   - basic alphanumeric/special chars string
	   - unicode string
	   - emoji string
	*/
	pub := newBasicPublishing("hello world")
	assert.Equal(suite.T(), string(pub.Body), "hello world")
}

func (suite *MQPTestSuite) TestPublishMessage() {
	/*
	   cases:
	   - basic known good single message
	   - multiple messages
	*/
	//conn := setup()
	conn := suite.mq.Connection
	channel, err := conn.Channel()
	assert.Nil(suite.T(), err)
	routingKey := "myRoutingKey"
	_, err = channel.QueueDeclare(
		routingKey, false, true, true, true, nil)
	assert.Nil(suite.T(), err)
	publishMessage(channel, "hello world")

	// got should be the last message in whatever test queue we used
	msg, _, err := channel.Get(routingKey, false)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), string(msg.Body), "hello world")
	channel.Close()
	//defer conn.Close()
}

func (suite *MQPTestSuite) TestProcessMessages() {
	//suite.T().Skip("pls fix me")
	/*
	   cases:
	   - one message
	   - no messages
	   - 5 messages
	*/
	conn := suite.mq.Connection
	channel, err := conn.Channel()
	assert.Nil(suite.T(), err)

	routingKey := "myRoutingKey"
	_, err = channel.QueueDeclare(
		routingKey, false, true, true, true, nil)
	assert.Nil(suite.T(), err)

	publishMessage(channel, "ping")

	messages := messages(channel, routingKey, "myConsumer")
	processMessages(messages)
}

func (suite *MQPTestSuite) TestParseFlags() {
	got := suite.mq.URI.String()
	want := "amqp://mqp:mqptest@127.0.0.1/mqp"
	assert.Equal(suite.T(), got, want)
}

/*
func TestConnectTLS(t *testing.T) {
	//got := fmt.Sprintf("%+v", conn)
	got := "foo"
	want := "bar"
	assert.Equal(t, got, want)
}
*/
