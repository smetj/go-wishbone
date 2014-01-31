package roundrobin

import "wishbone"
import "wishbone/event"
import "math/rand"

// import "fmt"

func NewModule(name string) actor.Actor {
	fanout := actor.NewActor()
	fanout.SetName(name)
	fanout.CreateQueue("inbox", 0)
	fanout.PreHook = PreHook
	return fanout
}

func PreHook(a *actor.Actor) {

	destination_list := []chan event.Event{}

	for index, _ := range a.Queuepool {
		if index != "inbox" && index != "_logs" && index != "_metrics" {
			destination_list = append(destination_list, a.Queuepool[index].Queue)
		}
	}

	c := generateConsumer(destination_list)
	a.RegisterConsumer(c, "inbox")
}

func generateConsumer(destination_list []chan event.Event) func(event.Event) {
	return func(e event.Event) {
		destination_list[random(0, len(destination_list))] <- e
	}
}

func random(min, max int) int {
	return rand.Intn(max-min) + min
}
