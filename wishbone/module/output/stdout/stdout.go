package stdout

import "wishbone"
import "wishbone/event"

import "fmt"

func NewModule(name string, silent bool) actor.Actor {
	stdout := actor.NewActor()
	stdout.SetName(name)

	stdout.CreateQueue("inbox", 0)
	c := generateConsumer(stdout.Name, silent)
	stdout.RegisterConsumer(c, "inbox")
	return stdout
}

func generateConsumer(name string, silent bool) func(event.Event) {
	if silent == true {
		return func(event event.Event) {
			//
		}
	} else {
		return func(event event.Event) {
			if event.Data != "" {
				fmt.Printf("%v - %v\n", name, event.Data)
			}
		}
	}
}
