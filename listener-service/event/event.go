package event

import amqp "github.com/rabbitmq/amqp091-go"

func declareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		"logs_topic", // name of the exchange
		"topic",      // type
		true,         // is exchange durable
		false,        //auto-deleted
		false,        // interna?
		false,        //no-wait?
		nil,          //arguments
	)

}

func declareRandomQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"RandomQueue", //name
		false,         // durable?
		false,         // delete when unused?
		true,          // is this exclusive?
		false,         // no wait?
		nil,           //any specific arguments
	)
}
