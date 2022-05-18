package event_bus

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
)

type EventBus struct {
	conn  *nats.Conn
	queue string
}

type subscriptionHandler func(topic string, event Msg) error

func NewEventBus(queue string) *EventBus {
	conn, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("connected to NATs")
	return &EventBus{conn: conn, queue: queue}
}

func (e *EventBus) Subscribe(topic string, handler subscriptionHandler) {
	if e.queue == "" {
		e.subscribe(topic, handler)
		return
	}
	e.queueSubscribe(topic, handler)
	return
}

func (e *EventBus) Publish(topic string, data interface{}) {
	b, _ := json.Marshal(data)
	fmt.Println(fmt.Sprintf("publishing to topic %s: %v", topic, data))
	h := nats.Header{}
	h.Set("service", e.queue)
	err := e.conn.PublishMsg(&nats.Msg{Header: h, Data: b, Subject: topic})
	if err != nil {
		fmt.Println("publish error: " + err.Error())
	}
}

func (e *EventBus) subscribe(topic string, handler func(string, Msg) error) {
	_, err := e.conn.Subscribe(e.queue, func(m *nats.Msg) {
		fmt.Println("message received on topic " + topic)
		err := handler(topic, Msg{Header: m.Header, Data: m.Data})
		if err != nil {
			return
		}
	})
	if err != nil {
		return
	}
	return
}

func (e *EventBus) queueSubscribe(topic string, handler func(string, Msg) error) {
	_, err := e.conn.QueueSubscribe(topic, e.queue, func(m *nats.Msg) {
		fmt.Println("message received on topic " + topic)
		err := handler(topic, Msg{Header: m.Header, Data: m.Data})
		if err != nil {
			return
		}
	})
	if err != nil {
		return
	}
	return
}
