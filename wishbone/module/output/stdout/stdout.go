package stdout

import "wishbone"
import "wishbone/event"
import "fmt"


func NewModule(name string)actor.Actor{
    stdout := actor.NewActor()
    stdout.SetName(name)

    stdout.CreateQueue("inbox")
    stdout.RegisterConsumer(consume, "inbox")
    return stdout
}

func consume(event event.Event){
    if event.Data != ""{
        fmt.Println(event.Data)
    }
}

