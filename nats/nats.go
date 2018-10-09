package nats

import (
	"bytes"
	"encoding/gob"

	"github.com/nats-io/go-nats"
)

type Nats struct {
	nc                      *nats.Conn
	logsCreatedSubscription *nats.Subscription
	logsCreatedChan         chan LogsCreatedMessage
}

func NewNats(url string) (*Nats, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return &Nats{nc: nc}, nil
}

func (e *Nats) SubscribeLogsCreated() (<-chan LogsCreatedMessage, error) {
	m := LogsCreatedMessage{}
	e.logsCreatedChan = make(chan LogsCreatedMessage, 64)
	ch := make(chan *nats.Msg, 64)
	var err error
	e.logsCreatedSubscription, err = e.nc.ChanSubscribe(m.Key(), ch)
	if err != nil {
		return nil, err
	}
	// Decode message
	go func() {
		for {
			select {
			case msg := <-ch:
				e.readMessage(msg.Data, &m)
				e.logsCreatedChan <- m
			}
		}
	}()
	return (<-chan LogsCreatedMessage)(e.logsCreatedChan), nil
}

func (e *Nats) OnLogsCreated(f func(LogsCreatedMessage)) (err error) {
	m := LogsCreatedMessage{}
	e.logsCreatedSubscription, err = e.nc.Subscribe(m.Key(), func(msg *nats.Msg) {
		e.readMessage(msg.Data, &m)
		f(m)
	})
	return
}

func (e *Nats) Close() {
	if e.nc != nil {
		e.nc.Close()
	}
	if e.logsCreatedSubscription != nil {
		e.logsCreatedSubscription.Unsubscribe()
	}
	close(e.logsCreatedChan)
}
func (e *Nats) PublishLogsCreated(mm LogsCreatedMessage) error {
	m := LogsCreatedMessage{mm.ID, mm.LogContent, mm.CreatedAt}
	data, err := e.writeMessage(&m)
	if err != nil {
		return err
	}
	return e.nc.Publish(m.Key(), data)
}

func (mq *Nats) writeMessage(m Message) ([]byte, error) {
	b := bytes.Buffer{}
	err := gob.NewEncoder(&b).Encode(m)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (mq *Nats) readMessage(data []byte, m interface{}) error {
	b := bytes.Buffer{}
	b.Write(data)
	//fmt.Println(b.String())
	return gob.NewDecoder(&b).Decode(m)
}
