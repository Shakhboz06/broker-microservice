package main

import (
	"context"
	"fmt"
	"listener-service/event"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	backoff "github.com/cenkalti/backoff/v4"
	amqp "github.com/rabbitmq/amqp091-go"
)

const rabbitMQURL = "amqp://guest:guest@rabbitmq"

func main() {

	// Creating a context that is cancelled on SIGINT/SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt, syscall.SIGTERM)
	defer stop()


	rabbitConn, err := amqpConnect(ctx, rabbitMQURL)
	if err != nil {
		log.Fatal(err)
	}

	defer rabbitConn.Close()

	log.Println("Connected to RabbitMQ messages ...")

	consumer, err := event.NewConsumer(rabbitConn)
	if err != nil{
		log.Panic(err)
		return 
	}

	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil{
		log.Println(err)
	}


	// Block until signal
	<-ctx.Done()
	log.Println("Shutting down")
}
 
func amqpConnect(ctx context.Context, url string) (*amqp.Connection, error) {

	// var backOff time.Duration = time.Second * 1
	// var counts int64
	// var connection *amqp.Connection

	exp := backoff.NewExponentialBackOff()
	exp.InitialInterval = 1 * time.Second
	exp.MaxInterval = 10 * time.Second
	exp.MaxElapsedTime = 2 * time.Minute

	bo := backoff.WithContext(exp, ctx)

	var conn *amqp.Connection

	operation := func() error {

		dialCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		
		var err error
		cfg := amqp.Config{
			Dial: func(network, addr string) (net.Conn, error) {
			  return (&net.Dialer{}).DialContext(dialCtx, network, addr)
			},
		}

		// underlying dial also gets a timeout
		conn, err = amqp.DialConfig(url, cfg )

		if err != nil {
			log.Printf("🔄 RabbitMQ dial failed: %v", err)
			return err
		}
		return nil
	}

	if err := backoff.Retry(operation, bo); err != nil {
		return nil, fmt.Errorf("rabbitmq connect retries exhausted: %w", err)
	}

	return conn, nil
}
