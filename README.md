go-wishbone
===========

A GO implementation of Python Wishbone

https://github.com/smetj/wishbone


Example
-------

    package main

    import "wishbone/router"
    import "wishbone/module/output/stdout"
    import "wishbone/module/input/testevent"
    import "wishbone/module/flow/funnel"



    func main() {

        router  := router.NewRouter()
        input1  := testevent.NewModule("input1", "Hello I am number one.")
        input2  := testevent.NewModule("input2", "Hello I am number two")
        funnel  := funnel.NewModule("funnel")
        output  := stdout.NewModule("output")

        router.Register(&input1)
        router.Register(&input2)
        router.Register(&funnel)
        router.Register(&output)

        router.Connect("input1.outbox", "funnel.input1")
        router.Connect("input2.outbox", "funnel.input2")
        router.Connect("funnel.outbox", "output.inbox")

        router.Start()
        router.Block()
    }
