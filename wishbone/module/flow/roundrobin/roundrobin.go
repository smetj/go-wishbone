package roundrobin

import "wishbone"
import "wishbone/event"



func NewModule(name string, outboxes []string)actor.Actor{
    roundrobin := actor.NewActor()
    roundrobin.SetName(name)
    roundrobin.CreateQueue("inbox")

    destinations :=[]chan event.Event{}

    for _, value := range outboxes {
        roundrobin.CreateQueue(value)
        destinations = append(destinations, roundrobin.GetQueue(value))
    }

    spread := generateSpreader(destinations)

    roundrobin.RegisterConsumer(spread, "inbox")

    return roundrobin
}

func generateSpreader(outboxes []chan event.Event)func(event.Event){
    return func(event event.Event){
        for _, value := range outboxes {
            value <- event
        }
    }
}

