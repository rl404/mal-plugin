package example

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rl404/mal-plugin/pubsub"
	"github.com/rl404/mal-plugin/pubsub/redis"
)

type sampleData struct {
	Field1 string
	Field2 int
}

// Producer is usage example of producer.
func Producer() {
	// Init client.
	client, err := redis.New("localhost:6379", "")
	if err != nil {
		panic(err)
	}

	// Don't forget to close.
	defer client.Close()

	// Sample data. Can be any type.
	data := sampleData{
		Field1: "a",
		Field2: 1,
	}

	// Publish data to specific topic/channel. Data will be encoded first.
	err = client.Publish("topic", data)
	if err != nil {
		panic(err)
	}
}

// Consumer is usage example of consumer.
func Consumer() {
	// Init client.
	client, err := redis.New("localhost:6379", "")
	if err != nil {
		panic(err)
	}

	// Don't forget to close.
	defer client.Close()

	// Subscribe to specific topic/channel.
	s, err := client.Subscribe("topic")
	if err != nil {
		panic(err)
	}

	// Need to convert to Channel interface first
	// because Go doesn't allow interface method to
	// return another interface.
	channel := s.(pubsub.Channel)

	// Don't forget to close subscription.
	defer channel.Close()

	// Prepare a new or existing variable for
	// incoming message. Data type should be the
	// same as when publish the message.
	var newData sampleData

	// Read incomming message. Message will be decoded
	// to newData. Don't forget to use pointer.
	msgs, errChan := channel.Read(&newData)

	// Prepare goroutine channel that will stop when
	// ctrl+c.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		// Loop for waiting incoming message.
		for {
			select {
			// If message comes.
			case <-msgs:
				// Proccess the message.
				fmt.Println(newData.Field1)

			// If error comes.
			case err = <-errChan:
				// Process the error.
				fmt.Println(err)
			}
		}
	}()

	<-sigChan
}
