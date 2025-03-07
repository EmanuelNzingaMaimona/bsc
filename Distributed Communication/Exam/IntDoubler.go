package main

import (
	"fmt"
	"log"
	"strconv"
	"github.com/streadway/amqp"
)

func main() {
	conn1, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn1.Close()

	conn2, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn2.Close()

	cho, err := conn1.Channel()
	failOnError(err, "Failed to open a channel")
	defer cho.Close()
	chi, err := conn2.Channel()
	failOnError(err, "Failed to open a channel")
	defer chi.Close()

	err = cho.ExchangeDeclare("p4Exchange", "direct", false, true, false, false, nil)
	failOnError(err, "Failed to declare an exchange")

	q, err := chi.QueueDeclare("intQueue", false, true, false, false, nil)
	failOnError(err, "Failed to declare a queue")

	err = chi.QueueBind(q.Name, "int", "p4Exchange", false, nil)
	failOnError(err, "Failed to bind a queue")

	msgs, err := chi.Consume(q.Name, "", false, false, false, false, nil)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			result := doubleIntegerString((string(d.Body)))
			err = cho.Publish(
				"p4Exchange", d.ReplyTo, false, false,
				amqp.Publishing{
					ContentType:   "text/plain",
					CorrelationId: d.CorrelationId,
					Body:          []byte(result),
				})
			failOnError(err, "Failed to publish a message")
			fmt.Println("Received job:", string(d.Body), "Published response:", result)
			d.Ack(false)
		}
	}()
	fmt.Println("Waiting for jobs")
	<-forever
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// Double integer strings "123" -> "246"
func doubleIntegerString(s string) string {
	var result string
	val, err := strconv.Atoi(string(s))
	if err == nil {
		result = strconv.Itoa((val * 2))
	}

	return result
}