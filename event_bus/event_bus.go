package event_bus

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"time"
)

type EventBus struct {
	conn  *nats.Conn
	queue string
}

type Event interface {
	SetTopic(string)
	SetData(interface{})
	SetOrigin(string)
	SetCorrelationId(string)
	GetTopic() string
	GetData() interface{}
	GetOrigin() string
	GetCorrelationId() string
}

type NATSEvent struct {
	topic         string
	data          interface{}
	origin        string
	correlationId string
}

func (n *NATSEvent) SetTopic(t string)         { n.topic = t }
func (n *NATSEvent) SetData(d interface{})     { n.data = d }
func (n *NATSEvent) SetOrigin(o string)        { n.origin = o }
func (n *NATSEvent) SetCorrelationId(c string) { n.correlationId = c }
func (n *NATSEvent) GetTopic() string          { return n.topic }
func (n *NATSEvent) GetData() interface{}      { return n.data }
func (n *NATSEvent) GetOrigin() string         { return n.origin }
func (n *NATSEvent) GetCorrelationId() string  { return n.correlationId }

type subscriptionHandler func(topic string, event Msg) error

func NewEventBus(queue string) *EventBus {
	conn, err := nats.Connect("nats://nats:4222")
	if err != nil {
		logError("initialising", "", err, "", "")
		panic(err)
	}
	fmt.Println("connected to NATs")
	return &EventBus{conn: conn, queue: queue}
}

func (e *EventBus) NewEvent() Event {
	return &NATSEvent{}
}

func (e *EventBus) Subscribe(topic string, handler subscriptionHandler) {
	if e.queue == "" {
		e.subscribe(topic, handler)
		return
	}
	e.queueSubscribe(topic, handler)
	return
}

func (e *EventBus) Publish(event Event) {
	b, err := json.Marshal(event.GetData())
	if err != nil {
		logError("publish error", event.GetTopic(), err, event.GetOrigin(), event.GetCorrelationId())
	}
	or := event.GetOrigin()
	cor := event.GetCorrelationId()
	h := nats.Header{}
	h.Set("service", e.queue)
	h.Set("origin", or)
	h.Set("correlation-id", cor)
	log("publish", event.GetTopic(), b, or, cor)
	err = e.conn.PublishMsg(&nats.Msg{Header: h, Data: b, Subject: event.GetTopic()})
	if err != nil {
		logError("publish error", event.GetTopic(), err, event.GetOrigin(), event.GetCorrelationId())
	}
}

func (e *EventBus) subscribe(topic string, handler func(string, Msg) error) {
	_, err := e.conn.Subscribe(e.queue, func(m *nats.Msg) {
		log("receive", topic, m.Data, m.Header.Get("origin"), m.Header.Get("correlation_id"))
		err := handler(topic, Msg{Header: m.Header, Data: m.Data})
		if err != nil {
			logError("receive error", topic, err, m.Header.Get("origin"), m.Header.Get("correlation_id"))
			return
		}
	})
	if err != nil {
		logError("subscription error", topic, err, "", "")
		return
	}
	return
}

func (e *EventBus) queueSubscribe(topic string, handler func(string, Msg) error) {
	_, err := e.conn.QueueSubscribe(topic, e.queue, func(m *nats.Msg) {
		log("queue receive", topic, m.Data, m.Header.Get("origin"), m.Header.Get("correlation_id"))
		err := handler(topic, Msg{Header: m.Header, Data: m.Data})
		if err != nil {
			logError("queue receive error", topic, err, m.Header.Get("origin"), m.Header.Get("correlation_id"))
			return
		}
	})
	if err != nil {
		logError("queue subscription error", topic, err, "", "")
		return
	}
	return
}

type logLine struct {
	Action        string      `json:"action"`
	CorrelationId string      `json:"correlation_id"`
	Origin        string      `json:"origin"`
	Topic         string      `json:"topic"`
	Data          interface{} `json:"data"`
	Time          time.Time   `json:"time"`
}

func log(action, topic string, data []byte, or, cor string) {
	l := logLine{
		Action:        action,
		Topic:         topic,
		Time:          time.Now(),
		Origin:        or,
		CorrelationId: cor,
	}
	var d interface{}
	_ = json.Unmarshal(data, &d)
	l.Data = d
	b, _ := json.Marshal(l)
	fmt.Println(string(b))
}

func logError(action, topic string, err error, or, cor string) {
	log(action, topic, []byte(err.Error()), or, cor)
}
