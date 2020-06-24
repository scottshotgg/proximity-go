package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/paulbellamy/ratecounter"
	grpc_consumer "github.com/scottshotgg/proximity-go/consumer/grpc"
	grpc_producer "github.com/scottshotgg/proximity-go/producer/grpc"
)

// func main() {
// 	const (
// 		nodeAddr = ":5001"
// 		route    = "a"
// 	)

// 	for i := 0; i < 1; i++ {
// 		go func() {
// 			var (
// 				contents []byte

// 				prod, err = grpc.New(nodeAddr)
// 			)

// 			if err != nil {
// 				log.Fatalln("err grpc.New:", err)
// 			}

// 			for {
// 				err = prod.Single(route, contents)
// 				if err != nil {
// 					log.Fatalln("err prod.Publish:", err)
// 				}
// 			}
// 		}()
// 	}

// 	time.Sleep(10 * time.Second)
// }
const (
	nodeAddr    = ":5001"
	route       = "a"
	everySecond = 1 * time.Second
	size        = 1
)

var (
	timer        = time.NewTimer(everySecond)
	sendcounter  = ratecounter.NewRateCounter(everySecond)
	recvcounter  = ratecounter.NewRateCounter(everySecond)
	sendcounters = make([]*ratecounter.RateCounter, size)
	recvcounters = make([]*ratecounter.RateCounter, size)
)

func main() {
	go func() {
		for {
			select {
			case <-timer.C:
				fmt.Printf("Totals:\nSender: %d\nRecver: %d\n", sendcounter.Rate(), recvcounter.Rate())

				for i, counter := range sendcounters {
					fmt.Printf("Sender %d: %d\n", i, counter.Rate())
				}

				for i, counter := range recvcounters {
					fmt.Printf("Recver %d: %d\n", i, counter.Rate())
				}

				fmt.Println()

				timer.Reset(everySecond)
			}
		}
	}()

	senders()
	recvers()

	time.Sleep(1000 * time.Second)
}

func recvers() {
	for i := 0; i < size; i++ {
		go func(i int) {
			recvcounters[i] = ratecounter.NewRateCounter(everySecond)
			time.Sleep(15 * time.Second)

			var (
				cons, err = grpc_consumer.New(nodeAddr, "id_"+strconv.Itoa(i), route+strconv.Itoa(i))
			)

			if err != nil {
				log.Fatalln("err grpc.New:", err)
			}

			var ch = make(chan []byte, 1000)

			cons.Stream(ch)

			for {
				<-ch

				recvcounters[i].Incr(1)
				recvcounter.Incr(1)
			}
		}(i)
	}
}

func senders() {
	for i := 0; i < size; i++ {
		go func(i int) {
			sendcounters[i] = ratecounter.NewRateCounter(everySecond)

			var (
				contents []byte

				prod, err = grpc_producer.New(nodeAddr)
			)

			if err != nil {
				log.Fatalln("err grpc.New:", err)
			}

			var ch = make(chan []byte, 100000)

			prod.Stream(route+strconv.Itoa(i), ch)

			var timer = time.After(5 * time.Second)

			for {
				select {
				case <-timer:
					return

				default:
				}

				ch <- contents

				sendcounters[i].Incr(1)
				sendcounter.Incr(1)
			}
		}(i)
	}
}