package main

import (
	"log"

	"github.com/nats-io/nats.go"
)

func main() {

	wait := make(chan bool)

	nc, err := nats.Connect(nats.DefaultURL)

	if err != nil {
		log.Fatal(err)
	}
	const topic = "send_voice_message"

	nc.Subscribe(topic, func(m *nats.Msg) {
		log.Printf("Received a message: %s\n", string(m.Data))

	})

	log.Println("Subscribed to", topic)

	<-wait
}
