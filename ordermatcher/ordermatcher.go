package ordermatcher

import (
	"container/heap"
	"fmt"
	"log"
	"main/kiteconnectsimulator/models"
	kiteticker "main/kiteconnectsimulator/ticker"

	"github.com/uptrace/bun"
)

var (
	buyOrderQueues        map[uint32]*PriorityQueue
	sellOrderQueues       map[uint32]*PriorityQueue
	subscribedInstruments map[uint32]bool

	BUY_OFFSET  = 0.10
	SELL_OFFSET = 0.10
)

type OrderMatcher struct {
	Ticker *kiteticker.Ticker
	Db     *bun.DB
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
		orderid:         orderid,
		quantity:        qty,
		price:           amount,
		transactionType: "BUY",
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
		orderid:         orderid,
		quantity:        qty,
		price:           amount,
		transactionType: "SELL",
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

func handleTick(tick models.Tick) {
	if queue, ok := buyOrderQueues[tick.InstrumentToken]; ok {
		order := queue.Top()
		if tick.LastPrice-order.price > BUY_OFFSET {
			// Complete_order_and_update_holding(order.orderid)
		}
	}
}
