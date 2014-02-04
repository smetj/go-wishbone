go-wishbone
===========

*A Go framework to build event pipeline servers with minimal effort.*

A Go implementation of Python Wishbone: https://github.com/smetj/wishbone


Example
-------

    package main

    import "wishbone/router"
    import "wishbone/module/output/stdout"
    import "wishbone/module/input/testevent"
    import "wishbone/module/flow/funnel"
    import "wishbone/module/system/metrics/graphite"
    import "wishbone/module/output/tcp"

    func main() {

        router := router.NewRouter()

        input1 := testevent.NewModule("input1", "Hello I am number one.")
        input2 := testevent.NewModule("input2", "Hello I am number two")
        funnel := funnel.NewModule("funnel")
        output := stdout.NewModule("output", false)

        log_output := stdout.NewModule("log_output", false)
        graphite := graphite.NewModule("graphite")
        metric_output := tcp.NewModule("metric_output", "graphite-001:2013", true, true)

        router.Register(&input1)
        router.Register(&input2)
        router.Register(&funnel)
        router.Register(&output)
        router.Register(&log_output)
        router.Register(&graphite)
        router.Register(&metric_output)

        // organize log flow
        router.Connect("_internal_logs.outbox", "log_output.inbox")

        // organise metric flow
        router.Connect("_internal_metrics.outbox", "graphite.inbox")
        router.Connect("graphite.outbox", "metric_output.inbox")

        // organize event flow
        router.Connect("input1.outbox", "funnel.input1")
        router.Connect("input2.outbox", "funnel.input2")
        router.Connect("funnel.outbox", "output.inbox")

        router.Start()
        router.Block()
    }

- Internal metrics are send to graphite.
- Logs are send to STDOUT


More Examples
-------------

More examples can be found at https://github.com/smetj/experiments/tree/master/go/go-wishbone