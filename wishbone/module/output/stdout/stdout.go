package stdout

import "wishbone"
import "wishbone/event"
import "fmt"


func NewModule(name string)actor.Actor{
    stdout := actor.NewActor()
    stdout.SetName(name)

    stdout.CreateQueue("inbox")
    c := generateConsumer(stdout.Name)
    stdout.RegisterConsumer(c, "inbox")
    return stdout
}

func generateConsumer(name string)func(event.Event){
    return func(event event.Event){
        if event.Data != ""{
            fmt.Printf("%v - %v\n",name, event.Data)
        }
    }
}

