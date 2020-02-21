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
	var err error
	suite.mq, err = NewMQ(uri)
	assert.Nil(suite.T(), err)
}

func (suite *MQPTestSuite) TearDownSuite() {
	suite.mq.Connection.Close()
}

func TestMQPTestSuite(t *testing.T) {
	suite.Run(t, new(MQPTestSuite))
}

func (suite *MQPTestSuite) TestConnect() {
	/*
		cases:
		- known good {credentials, vhost}
		- known bad {credentials, vhost}
		- rmq not running
	*/
	conn := suite.mq.Connection
	assert.Equal(suite.T(), "mqp", conn.Config.Vhost)
}

func (suite *MQPTestSuite) TestNewBasicPublishing() {
	/*
	   cases:
	   - basic alphanumeric/special chars string
	   - unicode string
	   - emoji string
	*/
	pub := newBasicPublishing("hello world")
	assert.Equal(suite.T(), "hello world", string(pub.Body))
}

func (suite *MQPTestSuite) TestPublishMessage() {
	/*
	   cases:
	   - basic known good single message
	   - multiple messages
	*/
	routingKey := "myRoutingKey"
	channel, err := suite.mq.setupChannel(routingKey)
	assert.Nil(suite.T(), err)
	publishMessage(channel, "hello world")

	// got should be the last message in whatever test queue we used
	msg, _, err := channel.Get(routingKey, false)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "hello world", string(msg.Body))
}

func (suite *MQPTestSuite) TestProcessMessages() {
	/*
	   cases:
	   - one message
	   - no messages
	   - 5 messages
	*/
	routingKey := "myRoutingKey"
	channel, err := suite.mq.setupChannel(routingKey)
	assert.Nil(suite.T(), err)

	publishMessage(channel, "ping")

	messages := messages(channel, routingKey, "myConsumer")
	//goroutine but close after X seconds
	go channelTimeout(channel, 1)
	processMessages(messages)
}

func (suite *MQPTestSuite) TestParseFlags() {
	got := suite.mq.URI.String()
	want := "amqp://mqp:mqptest@127.0.0.1/mqp"
	assert.Equal(suite.T(), want, got)
}

func (suite *MQPTestSuite) TestMain() {
	var err error
	assert.Nil(suite.T(), err)
}

/*
func TestConnectTLS(t *testing.T) {
	//got := fmt.Sprintf("%+v", conn)
	got := "foo"
	want := "bar"
	assert.Equal(t, got, want)
}
*/
