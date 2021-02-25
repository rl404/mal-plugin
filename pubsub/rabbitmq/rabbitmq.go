// Package rabbitmq is a wrapper of the original "github.com/streadway/amqp" library.
//
// Only contains basic publish, subscribe, and close methods.
// Data will be encoded to JSON before publishing the message.
package rabbitmq

import (
	"encoding/json"
	"sync/atomic"
	"time"

	"github.com/rl404/mal-plugin/pubsub"
	"github.com/streadway/amqp"
)

// Rabbitmq client implements PubSub interface.
var _ pubsub.PubSub = &Client{}

// Rabbitmq channels implements Channel interface.
var _ pubsub.Channel = &Channel{}

// Client is rabbitmq pubsub client.
type Client struct {
	client *amqp.Connection
}

// Channel is rabbitmq subscription channel.
type Channel struct {
	channel  *amqp.Channel
	messages <-chan amqp.Delivery
	closed   int32
}

const delay = 1

// New to create new rabbitmq pubsub client.
func New(url string) (*Client, error) {
	c, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	cl := &Client{client: c}

	// Auto reconnect.
	go func() {
		for {
			_, ok := <-cl.client.NotifyClose(make(chan *amqp.Error))
			if !ok {
				// Closed by function.
				break
			}
			for {
				time.Sleep(delay * time.Second)
				if cl.client, err = amqp.Dial(url); err == nil {
					// Reconnected.
					break
				}
			}
		}
	}()

	return cl, nil
}

// Publish to publish message.
func (c *Client) Publish(queue string, data interface{}) error {
	j, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if c.client == nil {
		return amqp.ErrClosed
	}

	ch, err := c.client.Channel()
	if err != nil {
		return err
	}

	q, err := ch.QueueDeclare(queue, false, false, false, false, nil)
	if err != nil {
		return err
	}

	err = ch.Publish("", q.Name, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        j,
	})
	if err != nil {
		return err
	}

	return ch.Close()
}

// Subscribe to subscribe queue.
func (c *Client) Subscribe(queue string) (interface{}, error) {
	ch, err := c.subscribe(queue)
	if err != nil {
		return nil, err
	}

	// Auto reconnect.
	go func() {
		for {
			_, ok := <-c.client.NotifyClose(make(chan *amqp.Error))
			if !ok || ch.isClosed() {
				// Closed by function.
				break
			}
			for {
				time.Sleep(delay * time.Second)
				if c.client == nil {
					// Wait until connection established.
					continue
				}
				tmp, err := c.subscribe(queue)
				if err == nil {
					// Reconnected.
					ch.messages = tmp.messages
					break
				}
			}
		}
	}()

	return ch, nil
}

func (c *Client) subscribe(queue string) (*Channel, error) {
	ch, err := c.client.Channel()
	if err != nil {
		return nil, err
	}

	_, err = ch.QueueDeclare(queue, false, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	msgs, err := ch.Consume(queue, "", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	return &Channel{
		channel:  ch,
		messages: msgs,
	}, nil
}

// Close to close pubsub connection.
func (c *Client) Close() error {
	return c.client.Close()
}

// Read to read incoming message.
func (c *Channel) Read(model interface{}) (<-chan interface{}, <-chan error) {
	msgChan, errChan := make(chan interface{}), make(chan error)
	go func() {
		for {
			for msg := range c.messages {
				if err := json.Unmarshal(msg.Body, &model); err != nil {
					errChan <- err
				} else {
					msgChan <- model
				}
			}

			time.Sleep(delay * time.Second)
			if c.isClosed() {
				break
			}
		}
	}()
	return (<-chan interface{})(msgChan), (<-chan error)(errChan)
}

// Close to close subscription.
func (c *Channel) Close() error {
	if c.isClosed() {
		return amqp.ErrClosed
	}
	atomic.StoreInt32(&c.closed, 1)
	return c.channel.Close()
}

func (c *Channel) isClosed() bool {
	return atomic.LoadInt32(&c.closed) == 1
}
