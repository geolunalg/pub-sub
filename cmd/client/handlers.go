package main

import (
	"fmt"

	"github.com/geolunalg/pub-sub/internal/gamelogic"
	"github.com/geolunalg/pub-sub/internal/routing"
)

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) {
	return func(ps routing.PlayingState) {
		defer fmt.Print("> ")
		gs.HandlePause(ps)
	}
}
