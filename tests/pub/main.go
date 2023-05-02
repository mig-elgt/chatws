package main

import (
	"bytes"
	"flag"
	"log"
	"strings"
	"time"

	"github.com/mig-elgt/broker"
)

func main() {
	topic := flag.String("topic", "foo:bar", "topic")
	topicMsg := flag.String("msg", "foo", "topic message")
	flag.Parse()
	// Create Event Queue object
	eq, err := broker.NewEventQueue("localhost", "5672", "user", "bitnami")
	if err != nil {
		panic(err)
	}
	defer eq.Close()
	// Open a new Channel
	ch, err := eq.OpenChannel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()
	t := strings.Split(*topic, ":")
	for {
		if err := eq.PublishEvent(t[0], t[1], bytes.NewBufferString(*topicMsg), ch); err != nil {
			log.Fatal(err)
		}
		time.Sleep(1 * time.Second)
	}
}
