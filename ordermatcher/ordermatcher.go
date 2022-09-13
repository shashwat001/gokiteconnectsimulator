package ordermatcher

import (
	"container/heap"
	"fmt"
	"log"
	kiteticker "main/kiteconnectsimulator/ticker"
)

var (
	buyOrderQueues        map[uint32]*PriorityQueue
	sellOrderQueues       map[uint32]*PriorityQueue
	subscribedInstruments map[uint32]bool
)

type OrderMatcher struct {
	Ticker *kiteticker.Ticker
}

func (om *OrderMatcher) Start() {
	buyOrderQueues = make(map[uint32]*PriorityQueue)
	sellOrderQueues = make(map[uint32]*PriorityQueue)
	subscribedInstruments = make(map[uint32]bool)

	log.Println("Starting listening to ticker")
	Runticker(om.Ticker)

}

func (om *OrderMatcher) AddBuy(orderid string, instrumentToken uint32, qty int64, amount float64) {
	log.Println("Adding buy order to limit queue for instrument: ", instrumentToken)
	ord := &order{
		orderid:  orderid,
		quantity: qty,
		price:    amount,
	}

	if queue, ok := buyOrderQueues[instrumentToken]; ok {
		heap.Push(queue, ord)
		fmt.Println(queue.Len())
	} else {
		queue := &PriorityQueue{}
		heap.Push(queue, ord)
		buyOrderQueues[instrumentToken] = queue
	}

	if _, ok := subscribedInstruments[instrumentToken]; !ok {
		log.Println("Subscribing to ticker for instrument: ", instrumentToken)
		subscribedInstruments[instrumentToken] = true
		om.Ticker.Subscribe([]uint32{instrumentToken})
	}
}

func (om *OrderMatcher) AddSell(orderid string, instrumentToken uint32, qty int64, amount float64) {
	log.Println("Adding sell order to limit queue for instrument: ", instrumentToken)
	ord := &order{
		quantity: qty,
		price:    amount,
	}

	if queue, ok := sellOrderQueues[instrumentToken]; ok {
		heap.Push(queue, ord)
		fmt.Println(queue.Len())
	} else {
		queue := &PriorityQueue{}
		heap.Push(queue, ord)
		sellOrderQueues[instrumentToken] = queue
	}

	if _, ok := subscribedInstruments[instrumentToken]; !ok {
		log.Println("Subscribing to ticker for instrument: ", instrumentToken)
		subscribedInstruments[instrumentToken] = true
		om.Ticker.Subscribe([]uint32{instrumentToken})
	}
}
