package actor

import "fmt"
import "time"
import "os"

// import "wishbone/logger"
import "wishbone/event"
import "code.google.com/p/go-uuid/uuid"

const (
	None  = 0
	Start = 1
	Stop  = 2
	Pause = 3
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

type Metric struct {
	Time   (int64)
	Type   (string)
	Source (string)
	Name   (string)
	Value  (uint64)
	Unit   (string)
	Tags   ([]string)
}

type Queue struct {
	Queue      chan (event.Event)
	admin      chan (int)
	function   interface{}
	state      int
	connected  bool
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
			if k != "_metrics" && k != "_logs" {
				a.Queuepool[k].admin <- Start
				a.Queuepool[k].state = Start
			}
		} else if a.Queuepool[k].connected == false {
			a.Log("debug", fmt.Sprintf("Queue %v is not connected and has no consumer registering nullConsumer.", k))
			a.RegisterConsumer(nullConsumer, k)
			a.Queuepool[k].admin <- Start
			a.Queuepool[k].state = Start
		}
	}
	a.Log("debug", "Start")
}

func (a *Actor) Stop(queue string) {
	if queue == "" {
		for k, _ := range a.Queuepool {
			if a.Queuepool[k].admin != nil {
				a.Queuepool[k].admin <- Stop
				a.Queuepool[k].state = Stop
			}
		}
		a.Log("debug", "Stopping all queues.")
	} else {
		a.Queuepool[queue].admin <- Stop
		a.Queuepool[queue].state = Stop
		a.Log("debug", fmt.Sprintf("Stopping queue %v.", queue))
	}
}

func (a *Actor) Pause() {
	for k, _ := range a.Queuepool {
		if a.Queuepool[k].admin != nil {
			a.Queuepool[k].admin <- Pause
			a.Queuepool[k].state = Pause
		}
	}
	a.Log("debug", "Pause")
}

func (a *Actor) CreateQueue(name string, size int) {
	var tmp = new(Queue)
	tmp.Queue = make(chan event.Event, size)
	tmp.state = None
	a.Queuepool[name] = tmp
	a.Log("debug", fmt.Sprintf("Queue %v created with size of %v", name, size))
}

func (a *Actor) GetQueue(name string) chan event.Event {
	return a.Queuepool[name].Queue
}
func (a *Actor) SetQueue(name string, q chan event.Event) {
	a.Queuepool[name].Queue = q
}

func (a *Actor) HasQueue(name string) bool {
	if _, ok := a.Queuepool[name]; ok {
		return true
	} else {
		return false
	}
}

func (a *Actor) MarkQueueConnected(name string) {
	a.Queuepool[name].connected = true
}
func (a *Actor) IsQueueConnected(name string) bool {
	return a.Queuepool[name].connected
}

func (a *Actor) RegisterConsumer(c func(event.Event) error, q string) {
	var tmp = a.Queuepool[q]
	tmp.function = c
	tmp.admin = make(chan int)
	tmp.total = 0
	tmp.prev_total = 0
	tmp.state = Pause
	a.Queuepool[q] = tmp
	go a.consumer(q)
}

func (a *Actor) HasConsumer(name string) bool {
	if a.Queuepool[name].state == None {
		return false
	} else {
		return true
	}
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
				err := a.Queuepool[queue].function.(func(event.Event) error)(e)

				if err == nil {
					a.Queuepool[queue].IncrementTotalHits()
				} else {
					a.Log("critical", fmt.Sprintf("An error occured.  Reason: %v", err))
					if _, ok := a.Queuepool["failed"]; ok {
						a.Queuepool["failed"].Queue <- e
						a.Queuepool["failed"].IncrementTotalHits()
					} else {
						a.Log("warning", "Module has no failed queue.  Dropping event.")
					}
				}
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
				a.Queuepool[queue].Queue <- e
				a.Queuepool[queue].IncrementTotalHits()
			}
		}
	}
	a.Log("debug", "Exit")
}

func (a *Actor) metricGatherer() {

	hostname, _ := os.Hostname()

	for {
		for queue, _ := range a.Queuepool {
			total := a.Queuepool[queue].total
			rate := total - a.Queuepool[queue].prev_total

			var tmp = a.Queuepool[queue]
			tmp.prev_total = total
			a.Queuepool[queue] = tmp

			metrics := []Metric{}
			metrics = append(metrics, Metric{Time: time.Now().Unix(), Type: "wishbone", Source: hostname, Name: fmt.Sprintf("%v.queue.%v.rate", a.Name, queue), Value: rate, Unit: "msg/s", Tags: nil})
			metrics = append(metrics, Metric{Time: time.Now().Unix(), Type: "wishbone", Source: hostname, Name: fmt.Sprintf("%v.queue.%v.size", a.Name, queue), Value: uint64(len(a.Queuepool[queue].Queue)), Unit: "msg", Tags: nil})

			for _, metric := range metrics {
				m := event.NewEvent()
				m.Data = metric
				a.Queuepool["_metrics"].Queue <- m
			}

		}
		time.Sleep(time.Second * 1)
	}
}

func nullConsumer(e event.Event) error {
	return nil
}
