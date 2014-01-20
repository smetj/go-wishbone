package testevent

import "wishbone"
import "wishbone/event"
// import "time"


func NewModule(name string, data string)actor.Actor{
    testevent:= actor.NewActor()
    testevent.SetName(name)

    generator := generateProduce(data)

    testevent.CreateQueue("outbox")
    testevent.RegisterProducer(generator, "outbox")
    return testevent
}

func generateProduce(data string)func()event.Event{
    return func()event.Event{
        e := event.NewEvent()
        e.Data=data
        return e
    }
}

