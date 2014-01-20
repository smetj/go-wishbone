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

func NewActor()Actor{
    a := Actor{}
    a.Queuepool = make(map[string]*Queue)
    a.logs = make(chan Log, 10)
    a.Name = uuid.New()
    return a
}

type Log struct{
    Level(string)
    Time(int64)
    Source(string)
    Message(string)
}

type Queue struct{
    Queue chan(event.Event)
    admin chan(int)
    function interface{}
    total uint64
    prev_total uint64
}

func (q *Queue) IncrementTotalHits(){
    q.total++
}

type Actor struct{
    Name string
    Queuepool map[string]*Queue
    logs chan(Log)
}

func (a *Actor) GetName()string{
    return a.Name
}

func (a *Actor) SetName(name string){
    a.Name = name
}

func (a *Actor) log(level string, message string){
    a.logs <- Log{Level:level, Time:time.Now().Unix(), Source:a.Name, Message: message }
}

func (a *Actor) GetLogChannel()chan(Log){
    return a.logs
}

func (a *Actor) Start(){
    for k, _ := range a.Queuepool{
        if a.Queuepool[k].admin != nil{
            a.Queuepool[k].admin <- Start
        }
    }
    go a.logGatherer()
    go a.metricGatherer()
    a.log("debug","Start")
}

func (a *Actor) Stop(){
    for k, _ := range a.Queuepool{
        if a.Queuepool[k].admin != nil{
            a.Queuepool[k].admin <- Stop
        }
    }
    a.log("debug","Stop")
}

func (a *Actor) Pause(){
    for k, _ := range a.Queuepool{
        if a.Queuepool[k].admin != nil{
            a.Queuepool[k].admin <- Pause
        }
    }
    a.log("debug","Pause")
}

func (a *Actor) CreateQueue(name string){
    var tmp = new(Queue)
    tmp.Queue = make(chan event.Event)
    a.Queuepool[name] = tmp
    a.log("debug",fmt.Sprintf("Queue %v created.",name))
}

func (a *Actor) GetQueue(name string)chan event.Event{
    if _, ok := a.Queuepool[name]; ok{
        //
    } else {
        a.CreateQueue(name)
    }
    return a.Queuepool[name].Queue
}

func (a *Actor) RegisterConsumer(c func(event.Event), q string){
    var tmp = a.Queuepool[q]
    tmp.function = c
    tmp.admin = make(chan int)
    tmp.total = 0
    tmp.prev_total = 0
    a.Queuepool[q] = tmp
    go a.consumer(q)
}

func (a *Actor) RegisterProducer(p func()event.Event, q string){
    var tmp = a.Queuepool[q]
    tmp.function = p
    tmp.total = 0
    tmp.prev_total = 0
    tmp.admin = make(chan int)
    a.Queuepool[q] = tmp
    go a.producer(q)
}

func (a *Actor) consumer(queue string){
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
                e := <- a.Queuepool[queue].Queue
                e.Header[a.Name].IncrementTTL()
                a.Queuepool[queue].function.(func(event.Event))(e)
                a.Queuepool[queue].IncrementTotalHits()
            }
        }
    }
    a.log("debug","Exit")
}

func (a *Actor) producer(queue string){
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
                e := a.Queuepool[queue].function.(func()event.Event)()
                e.Header[a.Name].IncrementTTL()
                a.Queuepool[queue].Queue <- e
                a.Queuepool[queue].IncrementTotalHits()
            }
        }
    }
    a.log("debug","Exit")
}

func (a *Actor) metricGatherer(){
    a.log("debug","Start metricGatherer")
    for {
        for queue, _ := range a.Queuepool {
            total := a.Queuepool[queue].total
            rate := total - a.Queuepool[queue].prev_total

            var tmp = a.Queuepool[queue]
            tmp.prev_total = total
            a.Queuepool[queue] = tmp

            a.log("info",fmt.Sprintf ("Module: %s, Queue: %s, Rate: %d", a.Name, queue , rate ) )
            time.Sleep(time.Second * 1)
        }
    }
}

func (a *Actor) logGatherer(){
    for {
        fmt.Println(<- a.logs)
    }
}