package ordermatcher

import (
	"container/heap"
	"fmt"
	"kiteconnectsimulator/db"
	"kiteconnectsimulator/models"
	kiteticker "kiteconnectsimulator/ticker"
	"log"
)

var (
	buyOrderQueues        map[uint32]*PriorityQueue
	sellOrderQueues       map[uint32]*PriorityQueue
	subscribedInstruments map[uint32]bool

	BUY_OFFSET  = 0.10
	SELL_OFFSET = 0.10
)

type OrderMatcher struct {
	Ticker          *kiteticker.MainTicker
	CallbacksTicker *kiteticker.Ticker
	Db              *db.DbClient
}

func (om *OrderMatcher) Start() {
	buyOrderQueues = make(map[uint32]*PriorityQueue)
	sellOrderQueues = make(map[uint32]*PriorityQueue)
	subscribedInstruments = make(map[uint32]bool)
	om.CallbacksTicker.OnSubscribe(om.onSubscribe)

	log.Println("Starting listening to ticker")
	Runticker(om)

}

func (om *OrderMatcher) AddBuy(orderid int64, instrumentToken uint32, qty int64, amount float64) {
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

func (om *OrderMatcher) AddSell(orderid int64, instrumentToken uint32, qty int64, amount float64) {
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

func (om *OrderMatcher) handleTick(tick models.Tick) {
	// log.Println(tick)
	go om.CallbacksTicker.TriggerTick(tick)

	hasBuy := false
	hasSell := false

	if buyqueue, ok := buyOrderQueues[tick.InstrumentToken]; ok {
		for buyqueue.Len() > 0 {
			order := buyqueue.Top()
			// log.Println("Order price: ", order.price)
			// log.Println("Tick price: ", tick.LastPrice)
			if order.price-tick.LastPrice > BUY_OFFSET {
				log.Println("Buy order complete")
				buyqueue.Pop()
				dbOrder := om.Db.Complete_order_and_update_holding(order.orderid)
				om.CallbacksTicker.TriggerOrderUpdate(dbOrder.Order)
			} else {
				break
			}
		}

		if buyqueue.Len() > 0 {
			hasBuy = true
		} else {
			delete(buyOrderQueues, tick.InstrumentToken)
		}
	}

	if sellqueue, ok := sellOrderQueues[tick.InstrumentToken]; ok {
		for sellqueue.Len() > 0 {
			order := sellqueue.Top()
			if tick.LastPrice-order.price > SELL_OFFSET {
				log.Println("Sell order complete")
				sellqueue.Pop()
				dbOrder := om.Db.Complete_order_and_update_holding(order.orderid)
				om.CallbacksTicker.TriggerOrderUpdate(dbOrder.Order)
			} else {
				break
			}
		}
		if sellqueue.Len() > 0 {
			hasSell = true
		} else {
			delete(sellOrderQueues, tick.InstrumentToken)
		}
	}

	if !(hasBuy || hasSell) {
		if _, ok := subscribedInstruments[tick.InstrumentToken]; ok {
			// log.Println("Unsubscribing to ticker for instrument: ", tick.InstrumentToken)
			// om.Ticker.Unsubscribe([]uint32{tick.InstrumentToken})
			// delete(subscribedInstruments, tick.InstrumentToken)
		} else {
			// log.Fatal("Tick received for unsubscribed token: ", tick.InstrumentToken)
		}
	}

}

func (om *OrderMatcher) onConnect() {
	om.CallbacksTicker.TriggerConnect()
}

func (om *OrderMatcher) onSubscribe(token []uint32) error {
	return om.Ticker.Subscribe(token)
}
