package router

import "wishbone/event"

import "wishbone/module/flow/funnel"
import "time"
import "strings"
import "fmt"

func NewRouter() Router {
	r := Router{}
	r.module = make(map[string]Module)
	// Register a Funnel module to collect all logs.
	f := funnel.NewModule("_logs")
	r.Register(&f)
	return r
}

type Router struct {
	module map[string]Module
}

type Module interface {
	Start()
	Stop()
	Pause()
	GetName() string
	GetQueue(name string) chan event.Event
}

func (r Router) Forward(source chan event.Event, destination chan event.Event) {
	go func() {
		for {
			event := <-source
			destination <- event
		}
	}()
}

func (r Router) Connect(source string, destination string) {
	s := strings.Split(source, ".")
	d := strings.Split(destination, ".")
	src := r.module[s[0]].GetQueue(s[1])
	dst := r.module[d[0]].GetQueue(d[1])
	go func() {
		for {
			event := <-src
			dst <- event
		}
	}()
}

func (r Router) Register(module Module) {
	name := module.GetName()
	r.module[name] = module
}

func (r Router) Start() {

	for _, v := range r.module {
		v.Start()
		r.Connect(fmt.Sprintf("%v.%v", v.GetName(), "_logs"), fmt.Sprintf("_logs.%v", v.GetName()))
	}
}

func (r Router) Pause() {
	for _, v := range r.module {
		v.Pause()
	}
}

func (r Router) Block() {
	for {
		time.Sleep(time.Second * 1)
	}
}
