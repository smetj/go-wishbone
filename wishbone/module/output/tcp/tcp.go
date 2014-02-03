package tcp

import "wishbone"

import "wishbone/event"
import "net"
import "os"
import "fmt"

// import "reflect"

func NewModule(name string, address string, success bool, failed bool) actor.Actor {
	tcp := actor.NewActor()
	tcp.SetName(name)

	tcp.CreateQueue("inbox", 0)
	tcp.CreateQueue("success", 0)
	tcp.CreateQueue("failed", 0)

	c := generateConsumer(address, tcp)
	tcp.RegisterConsumer(c, "inbox")

	return tcp
}

func generateConsumer(address string, a actor.Actor) func(event.Event) {

	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		fmt.Println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Println("Dial failed:", err.Error())
		os.Exit(1)
	}

	return func(e event.Event) {

		_, err = conn.Write([]byte(e.Data.(actor.Metric).Name))
		if err != nil {
			fmt.Println("Write to server failed:", err.Error())
			os.Exit(1)
		}
	}
}
