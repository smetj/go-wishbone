package human

import "wishbone"
import "wishbone/event"
import "fmt"


type Human struct{
    actor.Actor
}

func (s *Human)consume(event event.Event){

    s.Queuepool["outbox"].Queue <- event
}

func Init(name string)Human{
    human := Human{}
    human.Name=name
    human.Init("screen")
    human.CreateQueue("inbox")
    human.CreateQueue("outbox")
    human.RegisterConsumer(human.consume, "inbox")
    return human
}
