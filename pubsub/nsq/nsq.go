// Package nsq is a wrapper of the original "github.com/nsqio/go-nsq" library.
//
// Only contains basic publish, subscribe, and close methods.
// Data will be encoded to JSON before publishing the message.
package nsq

import (
	"encoding/json"

	"github.com/nsqio/go-nsq"
	"github.com/rl404/mal-plugin/pubsub"
)

// NSQ client implements PubSub interface.
var _ pubsub.PubSub = &Client{}

// NSQ channels implements Channel interface.
var _ pubsub.Channel = &Channel{}

// Client is NSQ pubsub client.
type Client struct {
	address string
	config  *nsq.Config
}

// Channel is NSQ subscription channel.
type Channel struct {
	consumer *nsq.Consumer
	messages chan *nsq.Message
}

// New to create new NSQ pubsub client.
func New(address string) (*Client, error) {
	return NewWithConfig(address, nsq.NewConfig())
}

// NewWithConfig to create new NSQ pubsub client with config.
func NewWithConfig(address string, cfg *nsq.Config) (*Client, error) {
	return &Client{
		address: address,
		config:  cfg,
	}, nil
}

// Publish to publish message.
func (c *Client) Publish(topic string, data interface{}) error {
	j, err := json.Marshal(data)
	if err != nil {
		return err
	}

	p, err := nsq.NewProducer(c.address, c.config)
	if err != nil {
		return err
	}
	defer p.Stop()

	return p.Publish(topic, j)
}

// Subscribe to subscribe to a topic.
func (c *Client) Subscribe(topic string) (interface{}, error) {
	cc, err := nsq.NewConsumer(topic, "channel", c.config)
	if err != nil {
		return nil, err
	}

	m := make(chan *nsq.Message)
	h := &msgHandler{messages: m}
	cc.AddHandler(h)

	err = cc.ConnectToNSQD(c.address)
	if err != nil {
		return nil, err
	}

	return &Channel{
		consumer: cc,
		messages: m,
	}, nil
}

// Close to close pubsub connection.
func (c *Client) Close() error {
	return nil
}

// Read to read incoming message.
func (c *Channel) Read(model interface{}) (<-chan interface{}, <-chan error) {
	msgChan, errChan := make(chan interface{}), make(chan error)
	go func() {
		for msg := range c.messages {
			if err := json.Unmarshal(msg.Body, &model); err != nil {
				errChan <- err
			} else {
				msgChan <- model
			}
		}
	}()
	return (<-chan interface{})(msgChan), (<-chan error)(errChan)
}

// Close to close subscription.
func (c *Channel) Close() error {
	c.consumer.Stop()
	return nil
}

type msgHandler struct {
	messages chan *nsq.Message
}

// HandleMessage to handle incoming message.
func (h *msgHandler) HandleMessage(m *nsq.Message) error {
	h.messages <- m
	return nil
}