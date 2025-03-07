package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/streadway/amqp"
)

func main() {
	//Connection
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	printErrorAndExit(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	//Channel
	ch, err := conn.Channel()
	printErrorAndExit(err, "Failed to open a channel")
	defer ch.Close()

	//Exchange
	err = ch.ExchangeDeclare(
		"jobExchange", // name
		"direct",      // type
		false,         // durable
		true,          // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	printErrorAndExit(err, "Failed to declare an exchange")

	//Decalre and bind queue
	q, err := ch.QueueDeclare(
		"jobQueue", // name,,empty string let server generate id
		false,      // durable
		true,       // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	printErrorAndExit(err, "Failed to declare a queue")
	err = ch.QueueBind(
		q.Name,        // queue name
		"jobkey",      // routing key
		"jobExchange", // exchange
		false,
		nil)
	printErrorAndExit(err, "Failed to bind a queue")

	//Consume
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer,empty string let server generate id
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	printErrorAndExit(err, "Failed to register a consumer")

	go func() {
		for d := range msgs {
			bodyString := string(d.Body)
			c, err := strconv.Atoi(bodyString)

			if err == nil {
				fmt.Println("Received:", c*2)
			} else {
				fmt.Println("Received:", bodyString)
			}
			d.Ack(false)
		}
	}()

	fmt.Println("Waiting for msgs")
	forever := make(chan bool)
	<-forever
}

func printErrorAndExit(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, ":", err)
	}
}
