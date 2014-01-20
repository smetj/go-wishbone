package funnel

import "wishbone"
import "wishbone/event"



func NewModule(name string, inboxes []string)actor.Actor{
    funnel := actor.NewActor()
    funnel.SetName(name)
    funnel.CreateQueue("outbox")

    for _, value := range inboxes {
        funnel.CreateQueue(value)
        c := generateConsumer(funnel.Queuepool["outbox"].Queue)
        funnel.RegisterConsumer(c, value)
    }

    return funnel
}

func generateConsumer(output chan event.Event)func(event.Event){
    return func(event event.Event){
        output <- event
    }
}