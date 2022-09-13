package ordermatcher

import (
	"container/heap"
	"fmt"
	"log"
	kiteticker "main/kiteconnectsimulator/ticker"
)

var (
	orderQueues           map[string]*PriorityQueue
	subscribedInstruments map[uint32]bool
)

type OrderMatcher struct {
	Ticker *kiteticker.Ticker
}

func (om *OrderMatcher) Start() {
	orderQueues = make(map[string]*PriorityQueue)
	subscribedInstruments = make(map[uint32]bool)

	log.Println("Starting listening to ticker")
	Runticker(om.Ticker)

}

func (om *OrderMatcher) AddBuy(instrument string, instrumentToken uint32, qty int64, amount float64) {
	log.Println("Adding buy order to limit queue for instrument: ", instrument)
	ord := &order{
		quantity: qty,
		price:    amount,
	}

	if queue, ok := orderQueues[instrument+":BUY"]; ok {
		heap.Push(queue, ord)
		fmt.Println(queue.Len())
	} else {
		queue := &PriorityQueue{}
		heap.Push(queue, ord)
		orderQueues[instrument+":BUY"] = queue
	}

	if _, ok := subscribedInstruments[instrumentToken]; !ok {
		log.Println("Subscribing to ticker for instrument: ", instrument)
		subscribedInstruments[instrumentToken] = true
		om.Ticker.Subscribe([]uint32{instrumentToken})
	}
}

func (om *OrderMatcher) AddSell(instrument string, instrumentToken uint32, qty int64, amount float64) {
	log.Println("Adding sell order to limit queue for instrument: ", instrument)
	ord := &order{
		quantity: qty,
		price:    amount,
	}

	if queue, ok := orderQueues[instrument+":SELL"]; ok {
		heap.Push(queue, ord)
		fmt.Println(queue.Len())
	} else {
		queue := &PriorityQueue{}
		heap.Push(queue, ord)
		orderQueues[instrument+":SELL"] = queue
	}
}

func (om *OrderMatcher) listAll(instrument string) {
	queue := orderQueues[instrument]
	queue.List()
}
