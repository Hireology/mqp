package main

import (
	"fmt"
	//"log"
	"testing"
)

// test BEHAVIOR not implementation

/*
func TestMain(t *testing.T) {
    got := "foo"
    want := "bar"
    if got != want {
        t.Errorf("got %s; want %s", got, want)
    }
}
*/

func TestConnect(t *testing.T) {
	connectionString := fmt.Sprintf("%s://%s:%s@%s:%d/%s",
		scheme, user, password, host, port, vhost)
	conn := connect(connectionString)
	got := fmt.Sprintf("%s", conn.Config.Vhost)
	//got := fmt.Sprintf("%+v", conn)
	want := vhost
	if got != want {
		t.Errorf("got %s; want %s", got, want)
	}
	defer conn.Close()
}

func TestNewBasicPublishing(t *testing.T) {
	pub := newBasicPublishing("hello world")
	got := string(pub.Body[:])
	want := "hello world"
	if got != want {
		t.Errorf("got %s; want %s", got, want)
	}
}

func TestPublishMessage(t *testing.T) {
	connectionString := fmt.Sprintf("%s://%s:%s@%s:%d/%s",
		scheme, user, password, host, port, vhost)
	conn := connect(connectionString)
	channel, err := conn.Channel()
	failOnError("channel.open", err)
	publishMessage(channel, "hello world")

	// got should be the last message in whatever test queue we used
	msg, _, err := channel.Get("myRoutingKey", false)
	failOnError("channel.get", err)
	got := string(msg.Body[:])
	want := "hello world"
	if got != want {
		t.Errorf("got %s; want %s", got, want)
	}
	defer conn.Close()
}

func TestConsumeMessage(t *testing.T) {
	// publish a message
	connectionString := fmt.Sprintf("%s://%s:%s@%s:%d/%s",
		scheme, user, password, host, port, vhost)
	conn := connect(connectionString)
	channel, err := conn.Channel()
	failOnError("channel.open", err)
	publishMessage(channel, "ping")

	// then consume
	// FIXME the channel never seems to close
	consumeMessage(channel, "myRoutingKey", "myConsumer")

	got := "foo"
	want := "bar"
	if got != want {
		t.Errorf("got %s; want %s", got, want)
	}
}

/*
func TestConnectTLS(t *testing.T) {
    got := "foo"
    want := "bar"
    if got != want {
        t.Errorf("got %s; want %s", got, want)
    }
}
*/
