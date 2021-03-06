package router

import "wishbone/event"

import "wishbone/module/flow/funnel"

// import "time"
import "strings"
import "fmt"
import "os"
import "os/signal"

func NewRouter() Router {
	r := Router{}
	r.module = make(map[string]Module)

	// Register a Funnel module to collect all logs.
	f := funnel.NewModule("_internal_logs")
	r.Register(&f)

	// Register a Funnel module to collect all metrics.
	l := funnel.NewModule("_internal_metrics")
	r.Register(&l)

	return r
}

type Router struct {
	module map[string]Module
}

type Module interface {
	Start()
	Stop(string)
	Pause()
	GetName() string
	GetQueue(name string) chan event.Event
	SetQueue(name string, q chan event.Event)
	HasQueue(string) bool
	CreateQueue(string, int)
	MarkQueueConnected(string)
}

func (r *Router) Forward(source chan event.Event, destination chan event.Event) {
	go func() {
		for {
			event := <-source
			destination <- event
		}
	}()
}

func (r *Router) Connect(source string, destination string) {
	s := strings.Split(source, ".")
	d := strings.Split(destination, ".")

	if r.module[s[0]].HasQueue(s[1]) == false {
		r.module[s[0]].CreateQueue(s[1], 10)
	}
	if r.module[d[0]].HasQueue(d[1]) == false {
		r.module[d[0]].CreateQueue(d[1], 10)
	}

	dst := r.module[d[0]].GetQueue(d[1])
	src := r.module[s[0]].GetQueue(s[1])
	r.module[s[0]].MarkQueueConnected(s[1])
	// r.module[s[0]].SetQueue(s[1], r.module[d[0]].GetQueue(d[1]))

	go func() {
		for {
			event := <-src
			dst <- event
		}
	}()
}

func (r *Router) Register(module Module) {
	name := module.GetName()
	r.module[name] = module
}

func (r *Router) GetModules() <-chan string {
	ch := make(chan string)
	go func() {
		for _, module := range r.module {
			ch <- module.GetName()
		}
		close(ch)
	}()
	return ch
}

func (r *Router) Start() {
	for _, module := range r.module {
		r.Connect(fmt.Sprintf("%v.%v", module.GetName(), "_logs"), fmt.Sprintf("_internal_logs.%v", module.GetName()))
		r.Connect(fmt.Sprintf("%v.%v", module.GetName(), "_metrics"), fmt.Sprintf("_internal_metrics.%v", module.GetName()))
	}
	for _, module := range r.module {
		module.Start()
	}
}

func (r *Router) Stop() {
	for module := range r.GetModules() {
		fmt.Println(module)
	}
}

func (r *Router) Pause() {
	for _, v := range r.module {
		v.Pause()
	}
}

func (r *Router) Block() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	r.Stop()
}
