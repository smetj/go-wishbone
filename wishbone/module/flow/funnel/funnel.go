package funnel

import "wishbone"
import "wishbone/event"

// import "fmt"

func NewModule(name string) actor.Actor {
	funnel := actor.NewActor()
	funnel.SetName(name)
	funnel.CreateQueue("outbox", 0)
	funnel.PreHook = PreHook
	return funnel
}

func PreHook(a *actor.Actor) {
	for index, _ := range a.Queuepool {
		c := generateConsumer(a.Queuepool["outbox"].Queue)
		a.RegisterConsumer(c, index)
	}
}

func generateConsumer(output chan event.Event) func(event.Event) {
	return func(event event.Event) {
		output <- event
	}
}
