package skeleton

import "wishbone"
import "wishbone/event"


type Skeleton struct{
    actor.Actor
}

func (s *Skeleton)consume(event event.Event){
    s.Queuepool["outbox"].Queue <- event
}

func Init(name string)Skeleton{
    skeleton := Skeleton{}
    skeleton.Name=name
    skeleton.Init("screen")
    skeleton.CreateQueue("inbox")
    skeleton.CreateQueue("outbox")
    skeleton.RegisterConsumer(skeleton.consume, "inbox")
    return skeleton
}
