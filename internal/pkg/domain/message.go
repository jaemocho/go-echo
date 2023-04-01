package domain

import (
	"fmt"
	"reflect"

	"github.com/Shopify/sarama"
)

type Message interface {
	Message() map[string]interface{}
}

type MessageSender interface {
	SendMessage(message Message) error
	Close() error
}

type KafkaMessage struct {
	NewMessage string
	Topic      string
}

func (k *KafkaMessage) Message() map[string]interface{} {
	m := make(map[string]interface{})

	rv := reflect.ValueOf(k)

	for i := 0; i < rv.NumField(); i++ {
		m[rv.Field(i).Type().Name()] = rv.Field(i).Interface()
	}

	return m
}

type KafkaMessageSender struct {
	Sender sarama.SyncProducer
}

func NewKafkaMessageSender() MessageSender {
	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	p, err := sarama.NewSyncProducer([]string{
		"kafka-01:9092",
		"kafka-02:9092",
		"kafka-03:9092"}, config)
	if err != nil {
		panic(err)
	}
	return &KafkaMessageSender{Sender: p}
}

func (k *KafkaMessageSender) SendMessage(message Message) error {

	topic := fmt.Sprintf("%v", message.Message()["Topic"])
	msg := fmt.Sprintf("%v", message.Message()["NewMessage"])

	partition, offset, err := k.Sender.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(msg),
	})
	if err != nil {
		return err
	}
	fmt.Printf("%d/%d\n", partition, offset)
	return nil
}

func (k *KafkaMessageSender) Close() error {
	err := k.Sender.Close()
	if err != nil {
		return err
	}
	return nil
}

func test() {
	ms := NewKafkaMessageSender()
	msg := &KafkaMessage{
		NewMessage: "test",
		Topic:      "test",
	}

	ms.SendMessage(msg)
}
