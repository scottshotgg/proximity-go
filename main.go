package main

import (
	"fmt"
	"log"
	"time"

	"github.com/paulbellamy/ratecounter"
	"github.com/pkg/profile"
	grpc_consumer "github.com/scottshotgg/proximity-go/consumer/grpc"
	grpc_producer "github.com/scottshotgg/proximity-go/producer/grpc"
	"github.com/scottshotgg/proximity/pkg/listener"
	channel_recv "github.com/scottshotgg/proximity/pkg/recv/channel"
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
	// nodeAddr = "192.168.1.227:5001"
	nodeAddr = ":5001"
	route    = channel_recv.RouteNoOp
	// route       = "a"
	everySecond = 1 * time.Second
	size        = 10
	sep         = "========================"
)

var (
	ids = make([]string, size)

	timer        = time.NewTimer(everySecond)
	sendcounters = make([]*ratecounter.RateCounter, size)
	recvcounters = make([]*ratecounter.RateCounter, size)
)

func main() {
	defer profile.Start(profile.MemProfile, profile.ProfilePath("profile.p")).Stop()

	go func() {
		for {
			select {
			case <-timer.C:
				var (
					sendRate int64
					recvRate int64

					sendRates = make([]int64, size)
					recvRates = make([]int64, size)
				)

				for i := 0; i < size; i++ {
					sendRates[i] = sendcounters[i].Rate()

					sendRate += sendRates[i]
				}

				for i := 0; i < size; i++ {
					recvRates[i] = recvcounters[i].Rate()

					recvRate += recvRates[i]
				}

				fmt.Println("Totals:")
				fmt.Println(sep)

				fmt.Printf("\tSender: %d\n", sendRate)

				for i, j := range sendRates {
					fmt.Printf("\t%6d: %d\n", i, j)
				}

				fmt.Println(sep)

				fmt.Printf("\tRecver: %d\n", recvRate)

				for i, j := range recvRates {
					fmt.Printf("\t%6d: %d\n", i, j)
				}

				// fmt.Println(sep)
				fmt.Println()

				timer.Reset(everySecond)
			}
		}
	}()

	recvers()
	senders()

	time.Sleep(10 * time.Second)
}

func recvers() {
	for i := 0; i < size; i++ {
		go func(i int) {
			recvcounters[i] = ratecounter.NewRateCounter(everySecond)
			// time.Sleep(15 * time.Second)

			// var cons, err = grpc_consumer.New(nodeAddr, []string{route + strconv.Itoa(i)})
			var cons, err = grpc_consumer.New(nodeAddr, []string{route})
			if err != nil {
				log.Fatalln("err grpc.New:", err)
			}

			ids[i] = cons.ID()

			var ch = make(chan []byte, 1000)

			cons.Listen(ch)

			for {
				<-ch

				recvcounters[i].Incr(1)
			}
		}(i)
	}
}

func senders() {
	for i := 0; i < size; i++ {
		go func(i int) {
			sendcounters[i] = ratecounter.NewRateCounter(everySecond)
			// time.Sleep(2 * time.Second)

			var (
				contents []byte

				prod, err = grpc_producer.New(nodeAddr)
			)

			if err != nil {
				log.Fatalln("err grpc.New:", err)
			}

			var ch = make(chan *listener.Msg, 1000)

			prod.Stream(ch)

			// var timer = time.After(5 * time.Second)

			for {
				// select {
				// case <-timer:
				// 	return

				// default:
				// }

				if ids[i] == "" {
					continue
				}

				ch <- &listener.Msg{
					Route: channel_recv.RouteID + "/" + ids[i],
					// Route: route + strconv.Itoa(i),
					// Route:    route,
					Contents: contents,
				}

				// time.Sleep(20 * time.Millisecond)

				sendcounters[i].Incr(1)
			}
		}(i)
	}
}
