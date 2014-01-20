package logger

import "time"

type Log struct{
    Level(int)
    Time(int64)
    Source(string)
    Message(string)
}

type Logger interface{
    Alert(string)
    Critical(string)
    Error(string)
    Warning(string)
    Notice(string)
    Info(string)
    Debug(string)
}

type QueueLogger struct{
    Name string
    Logs chan(Log)
}


func (l QueueLogger) Alert(message string){
    l.Logs <- Log{Level:1, Time:time.Now().Unix(), Source:l.Name, Message: message }
}
func (l QueueLogger) Critical(message string){
    l.Logs <- Log{Level:2, Time:time.Now().Unix(), Source:l.Name, Message: message }
}
func (l QueueLogger) Error(message string){
    l.Logs <- Log{Level:3, Time:time.Now().Unix(), Source:l.Name, Message: message }
}
func (l QueueLogger) Warning(message string){
    l.Logs <- Log{Level:4, Time:time.Now().Unix(), Source:l.Name, Message: message }
}
func (l QueueLogger) Notice(message string){
    l.Logs <- Log{Level:5, Time:time.Now().Unix(), Source:l.Name, Message: message }
}
func (l QueueLogger) Info(message string){
    l.Logs <- Log{Level:6, Time:time.Now().Unix(), Source:l.Name, Message: message }
}
func (l QueueLogger) Debug(message string){
    l.Logs <- Log{Level:7, Time:time.Now().Unix(), Source:l.Name, Message: message }
}