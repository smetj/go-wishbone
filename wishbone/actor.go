package actor

import "fmt"
import "time"

// import "wishbone/logger"
import "wishbone/event"
import "code.google.com/p/go-uuid/uuid"

const (
	Start = 0
	Stop  = 1
	Pause = 2
)

func NewActor() Actor {
	a := Actor{}
	a.Queuepool = make(map[string]*Queue)
	a.Name = uuid.New()
	a.CreateQueue("_logs", 100)
	a.CreateQueue("_metrics", 100)
	return a
}

type Log struct {
	Level   (string)
	Time    (int64)
	Source  (string)
	Message (string)
}

type Queue struct {
	Queue      chan (event.Event)
	admin      chan (int)
	function   interface{}
	total      uint64
	prev_total uint64
}

func (q *Queue) IncrementTotalHits() {
	q.total++
}

type Actor struct {
	Name      string
	Queuepool map[string]*Queue
	PreHook   func(*Actor)
	PostHook  func(*Actor)
}

func (a *Actor) GetName() string {
	return a.Name
}

func (a *Actor) SetName(name string) {
	a.Name = name
}

func (a *Actor) Log(level string, message string) {
	l := event.NewEvent()
	l.Data = Log{Level: level, Time: time.Now().Unix(), Source: a.Name, Message: message}
	a.Queuepool["_logs"].Queue <- l
}

func (a *Actor) Start() {
	go a.metricGatherer()
	if a.PreHook != nil {
		a.Log("debug", "PreHook() found thus executing.")
		a.PreHook(a)
	} else {
		a.Log("debug", "No PreHook() found.")
	}

	for k, _ := range a.Queuepool {
		if a.Queuepool[k].admin != nil {
			a.Queuepool[k].admin <- Start
		}
	}

	a.Log("debug", "Start")
}

func (a *Actor) Stop() {
	for k, _ := range a.Queuepool {
		if a.Queuepool[k].admin != nil {
			a.Queuepool[k].admin <- Stop
		}
	}
	a.Log("debug", "Stop")
}

func (a *Actor) Pause() {
	for k, _ := range a.Queuepool {
		if a.Queuepool[k].admin != nil {
			a.Queuepool[k].admin <- Pause
		}
	}
	a.Log("debug", "Pause")
}

func (a *Actor) CreateQueue(name string, size int) {
	var tmp = new(Queue)
	tmp.Queue = make(chan event.Event, size)
	a.Queuepool[name] = tmp
	a.Log("debug", fmt.Sprintf("Queue %v created with size of %v", name, size))
}

func (a *Actor) GetQueue(name string) chan event.Event {
	return a.Queuepool[name].Queue
}

func (a *Actor) HasQueue(name string) bool {
	if _, ok := a.Queuepool[name]; ok {
		return true
	} else {
		return false
	}
}

func (a *Actor) RegisterConsumer(c func(event.Event), q string) {
	var tmp = a.Queuepool[q]
	tmp.function = c
	tmp.admin = make(chan int)
	tmp.total = 0
	tmp.prev_total = 0
	a.Queuepool[q] = tmp
	go a.consumer(q)
}

func (a *Actor) RegisterProducer(p func() event.Event, q string) {
	var tmp = a.Queuepool[q]
	tmp.function = p
	tmp.total = 0
	tmp.prev_total = 0
	tmp.admin = make(chan int)
	a.Queuepool[q] = tmp
	go a.producer(q)
}

func (a *Actor) consumer(queue string) {
	state := Pause
begin:
	for {
		select {
		case state = <-a.Queuepool[queue].admin:
			switch state {
			case Stop:
				break begin
			case Start:
				//
			case Pause:
				//
			}
		default:
			if state != Pause {
				e := <-a.Queuepool[queue].Queue
				e.Header[a.Name].IncrementTTL()
				a.Queuepool[queue].function.(func(event.Event))(e)
				a.Queuepool[queue].IncrementTotalHits()
			}
		}
	}
	a.Log("debug", "Exit")
}

func (a *Actor) producer(queue string) {
	state := Pause
begin:
	for {
		select {
		case state = <-a.Queuepool[queue].admin:
			switch state {
			case Stop:
				break begin
			case Start:
				//
			case Pause:
				//
			}
		default:
			if state != Pause {
				e := a.Queuepool[queue].function.(func() event.Event)()
				e.Header[a.Name].IncrementTTL()
				a.Queuepool[queue].Queue <- e
				a.Queuepool[queue].IncrementTotalHits()
			}
		}
	}
	a.Log("debug", "Exit")
}

func (a *Actor) metricGatherer() {

	for {
		for queue, _ := range a.Queuepool {
			total := a.Queuepool[queue].total
			rate := total - a.Queuepool[queue].prev_total

			var tmp = a.Queuepool[queue]
			tmp.prev_total = total
			a.Queuepool[queue] = tmp
			a.Log("info", fmt.Sprintf("Queue: %s, Rate: %d", queue, rate))
		}
		time.Sleep(time.Second * 1)
	}
}
