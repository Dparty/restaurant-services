package main

import (
	"time"

	"github.com/Dparty/restaurant-services/pubsub"
)

var pb = pubsub.GetPubSub()

func main() {
	c := pb.Subscribe("1234")
	c2 := pb.Subscribe("1234")
	defer c.Close()
	go func() {
		for {
			time.Sleep(1 * time.Second)
			pb.Publish("1234", "hello")
		}
	}()
	ch := c.Channel()
	ch2 := c2.Channel()
	for {
		select {
		case msg := <-ch:
			println("1:", msg.Payload)
		case msg := <-ch2:
			println("2:", msg.Payload)
		}
	}

}
