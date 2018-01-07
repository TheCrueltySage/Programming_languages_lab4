package main

import ("fmt";"flag")

type Token struct {
    data string
    recipient int
    ttl int
}

func ringMember(address int, input <-chan Token, output chan<- Token) {
    token := <-input
    token.ttl -= 1
    if token.recipient == address {
	fmt.Println("address:",address)
	fmt.Println("data:",token.data)
	output <- Token{"Die for me", -1, -1}
    } else if token.recipient <=0 {
	output<- Token{"Die for me", -1, -1}
    } else if token.ttl <=0 {
	fmt.Println("Time-to-live of message ran out on recpieint", address, ", dropping")
	output<- Token{"Die for me", -1, -1}
    } else {
	output<- token
    }
}

func main() {
    threads := flag.Int("members", 6, "amount of ring members spawned")
    data := flag.String("data", "", "data to send to ring member")
    rec := flag.Int("rec", 5, "receiver of the message")
    ttl := flag.Int("ttl", 6, "time-to-live of the message")
    flag.Parse()

    bus := make([]chan Token, *threads)
    for i:=0;i<*threads;i++ {
	bus[i] = make(chan Token)
    }

    status := Token{*data, *rec, *ttl}
    finished := false
    for finished == false {
	for i:=0;i<*threads-1;i++ {
	    go ringMember(i, bus[i], bus[i+1])
	}
	go ringMember(*threads-1,bus[*threads-1],bus[0])

	bus[0]<- status
	status = <-bus[0]
	if status.recipient<=0 {
	    finished = true
	}
    }
}
