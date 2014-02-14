package graphite

import "wishbone"
import "wishbone/event"
import "fmt"
import "os"

func NewModule(name string) actor.Actor {
	graphite := actor.NewActor()
	graphite.SetName(name)
	graphite.CreateQueue("inbox", 0)
	graphite.CreateQueue("outbox", 0)

	c := generateConsumer(graphite)
	graphite.RegisterConsumer(c, "inbox")

	return graphite
}

func generateConsumer(a actor.Actor) func(event.Event) error {

	return func(e event.Event) error {
		metric := fmt.Sprintf("%v.%v.%v %v %v\n", e.Data.(actor.Metric).Source, os.Args[0], e.Data.(actor.Metric).Name, e.Data.(actor.Metric).Value, e.Data.(actor.Metric).Time)
		e.Data = metric
		a.Queuepool["outbox"].Queue <- e
		return nil
	}
}
