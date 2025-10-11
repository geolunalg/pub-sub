package main

import (
	"fmt"
	"log"

	"github.com/geolunalg/pub-sub/internal/gamelogic"
	"github.com/geolunalg/pub-sub/internal/pubsub"
	"github.com/geolunalg/pub-sub/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	const rabbitConnString = "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(rabbitConnString)
	if err != nil {
		log.Fatalf("could not connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	fmt.Println("Peril game server connected to RabbitMQ!")

	publishCh, err := conn.Channel()
	if err != nil {
		log.Fatalf("could not subscribe to pause: %v", err)
	}

	_, queue, err := pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilDirect,
		routing.GameLogSlug,
		routing.GameLogSlug+".",
		pubsub.SimpleQueueDurable,
	)
	if err != nil {
		log.Fatalf("could not subscribe to pause: %v", err)
	}
	fmt.Printf("Queue %v declared and bound!\n", queue.Name)

	gamelogic.PrintServerHelp()

OuterLoop:
	for {
		words := gamelogic.GetInput()
		switch words[0] {
		case "pause":
			fmt.Println("sending a pause message")
			err = pubsub.PublishJSON(
				publishCh,
				routing.ExchangePerilDirect,
				routing.PauseKey,
				routing.PlayingState{
					IsPaused: true,
				},
			)
			if err != nil {
				log.Printf("could not publish time: %v", err)
			}
			fmt.Println("Pause message sent!")
		case "resume":
			fmt.Println("sending a resume message")
			err = pubsub.PublishJSON(
				publishCh,
				routing.ExchangePerilDirect,
				routing.PauseKey,
				routing.PlayingState{
					IsPaused: false,
				},
			)
			if err != nil {
				log.Printf("could not publish time: %v", err)
			}
			fmt.Println("Resume message sent!")
		case "quit":
			fmt.Println("exiting the program")
			break OuterLoop
		case "help":
			gamelogic.PrintServerHelp()
		default:
			fmt.Println("command unrecognized")
			gamelogic.PrintServerHelp()
		}
	}
}
