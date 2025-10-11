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
	fmt.Println("Starting Peril client...")
	const rabbitConnString = "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(rabbitConnString)
	if err != nil {
		log.Fatalf("could not connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	fmt.Println("Peril game client connected to RabbitMQ!")

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("could not get username: %v", err)
	}

	_, queue, err := pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilDirect,
		routing.PauseKey+"."+username,
		routing.PauseKey,
		pubsub.SimpleQueueTransient,
	)
	if err != nil {
		log.Fatalf("could not subscribe to pause: %v", err)
	}
	fmt.Printf("Queue %v declared and bound!\n", queue.Name)

	gamestate := gamelogic.NewGameState(username)

	err = pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilDirect,
		routing.PauseKey+"."+gamestate.GetUsername(),
		routing.PauseKey,
		pubsub.SimpleQueueTransient,
		handlerPause(gamestate),
	)
	if err != nil {
		log.Fatalf("could not subscribe to pause: %v", err)
	}

OuterLoop:
	for {
		words := gamelogic.GetInput()
		switch words[0] {
		case "spawn":
			err = gamestate.CommandSpawn(words)
			if err != nil {
				fmt.Println(err)
			}
		case "move":
			_, err := gamestate.CommandMove(words)
			if err != nil {
				fmt.Println(err)
			}
		case "status":
			gamestate.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			fmt.Println("Spamming not allowed yet!")
		case "quit":
			gamelogic.PrintQuit()
			break OuterLoop
		default:
			fmt.Println("Unknow command")
			gamelogic.PrintClientHelp()
		}
	}

	fmt.Println("RabbitMQ connection closed.")
}
