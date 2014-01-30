package fanout

import "wishbone"
import "wishbone/event"

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
		for _, destination := range destination_list {
			destination <- e
		}
	}
}
