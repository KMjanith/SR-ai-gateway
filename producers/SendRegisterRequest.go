package producers

import (
	"api-gateway/spec"
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"
)

func SendRegisterRequest(userName string, password string, ctx context.Context) (res []byte) {
	msg := &spec.RegisterUser{
		Username: userName,
		Password: password,
	}

	ch := ctx.Value("producerChannel")

	// Declare a queue
	q, err := ch.(*amqp.Channel).QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	if err != nil {
		log.Printf("cant declare the queu : %s", err)
	}

	// Consume the reply queue
	msgs, err := ch.(*amqp.Channel).Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	if err != nil {
		log.Printf("cant declare the queu : %s", err)
	}

	// Serialize message using protobuf
	request, err := proto.Marshal(msg)
	if err != nil {
		log.Fatal(err)
	}

	ctxx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.(*amqp.Channel).PublishWithContext(ctxx,
		"",     // exchange
		"ApiAuthMsg", // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/protobuf",
			ReplyTo:     q.Name,
			Body:        request,
		})
	if err != nil {
		log.Fatalf("Failed to publish a message: %v", err)
	}
	log.Printf(" [x] Sent %s\n", msg)

	for d := range msgs {

		log.Printf("Received response from sertingservice, tokem" )
		res = d.Body
		break

	}
	return
}
