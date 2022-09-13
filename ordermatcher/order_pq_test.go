package ordermatcher

import (
	"container/heap"
	"testing"

	"github.com/stretchr/testify/assert"
)

// test function
func TestBuyPQFunctions(t *testing.T) {
	queue := PriorityQueue{}
	ord := &order{orderid: 1, quantity: 10, price: 100.0, transactionType: "BUY"}
	ord1 := &order{orderid: 2, quantity: 100, price: 10.0, transactionType: "BUY"}
	ord2 := &order{orderid: 3, quantity: 15, price: 15.0, transactionType: "BUY"}
	ord3 := &order{orderid: 4, quantity: 110, price: 110.0, transactionType: "BUY"}
	ord4 := &order{orderid: 5, quantity: 110, price: 80.0, transactionType: "BUY"}
	ord5 := &order{orderid: 6, quantity: 110, price: 5.0, transactionType: "BUY"}
	heap.Push(&queue, ord)
	heap.Push(&queue, ord1)
	heap.Push(&queue, ord2)
	heap.Push(&queue, ord3)

	assert.Equal(t, 110.0, queue.Top().price)

	or := heap.Pop(&queue).(*order)
	assert.Equal(t, 110.0, or.price)
	assert.Equal(t, 100.0, queue.Top().price)

	or = heap.Pop(&queue).(*order)
	assert.Equal(t, 100.0, or.price)
	assert.Equal(t, 15.0, queue.Top().price)

	or = heap.Pop(&queue).(*order)
	assert.Equal(t, 15.0, or.price)
	assert.Equal(t, 10.0, queue.Top().price)

	heap.Push(&queue, ord4)
	assert.Equal(t, 80.0, queue.Top().price)
	or = heap.Pop(&queue).(*order)
	assert.Equal(t, 80.0, or.price)

	heap.Push(&queue, ord5)
	or = heap.Pop(&queue).(*order)
	assert.Equal(t, 10.0, or.price)

	assert.Equal(t, 5.0, queue.Top().price)
	or = heap.Pop(&queue).(*order)
	assert.Equal(t, 5.0, or.price)

	assert.Equal(t, 0, queue.Len())

}

func TestSellPQFunctions(t *testing.T) {
	queue := PriorityQueue{}
	ord := &order{orderid: 1, quantity: 10, price: 100.0, transactionType: "SELL"}
	ord1 := &order{orderid: 2, quantity: 100, price: 10.0, transactionType: "SELL"}
	ord2 := &order{orderid: 3, quantity: 15, price: 15.0, transactionType: "SELL"}
	ord3 := &order{orderid: 4, quantity: 110, price: 110.0, transactionType: "SELL"}
	ord4 := &order{orderid: 5, quantity: 110, price: 80.0, transactionType: "SELL"}
	ord5 := &order{orderid: 6, quantity: 110, price: 5.0, transactionType: "SELL"}
	heap.Push(&queue, ord)
	heap.Push(&queue, ord1)
	heap.Push(&queue, ord2)
	heap.Push(&queue, ord3)

	assert.Equal(t, 10.0, queue.Top().price)

	or := heap.Pop(&queue).(*order)
	assert.Equal(t, 10.0, or.price)
	assert.Equal(t, 15.0, queue.Top().price)

	or = heap.Pop(&queue).(*order)
	assert.Equal(t, 15.0, or.price)
	assert.Equal(t, 100.0, queue.Top().price)

	or = heap.Pop(&queue).(*order)
	assert.Equal(t, 100.0, or.price)
	assert.Equal(t, 110.0, queue.Top().price)

	heap.Push(&queue, ord4)
	assert.Equal(t, 80.0, queue.Top().price)
	or = heap.Pop(&queue).(*order)
	assert.Equal(t, 80.0, or.price)

	heap.Push(&queue, ord5)
	or = heap.Pop(&queue).(*order)
	assert.Equal(t, 5.0, or.price)

	assert.Equal(t, 110.0, queue.Top().price)
	or = heap.Pop(&queue).(*order)
	assert.Equal(t, 110.0, or.price)

	assert.Equal(t, 0, queue.Len())

}
