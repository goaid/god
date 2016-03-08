package god

import (
	"fmt"

	"github.com/streadway/amqp"
)

type Session struct {
	*amqp.Channel
}

func NewSession() (*Session, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	var s Session
	s.Channel = ch
	return &s, nil
}

func combine(routingKeyType uint16, routingKey uint64) string {
	return fmt.Sprintf("%d-%d", routingKeyType, routingKey)
}

func (s *Session) Declare(exchange string) error {
	return s.ExchangeDeclare(
		exchange, // name
		"direct", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
}

func (s *Session) Post(exchange string,
	routingKeyType uint16, routingKey uint64,
	msgID uint64, msg []byte) error {

	return s.Publish(
		exchange, // exchange
		combine(routingKeyType, routingKey), // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         msg,
		})
}

func (s *Session) Subscribe(exchange string,
	routingKeyType uint16, routingKey uint64) (string, error) {
	err := s.Declare(exchange)
	if err != nil {
		return "", err
	}

	q, err := s.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when usused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return "", err
	}

	err = s.QueueBind(
		q.Name, // queue name
		combine(routingKeyType, routingKey), // routing key
		exchange, // exchange
		false,
		nil)
	if err != nil {
		return "", err
	}

	return q.Name, nil

}

func (s *Session) Pull(queue string) (<-chan amqp.Delivery, error) {
	return s.Consume(
		queue, // queue
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
}
