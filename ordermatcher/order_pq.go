package ordermatcher

import (
	"container/heap"
	"fmt"
)

type order struct {
	orderid         string
	quantity        int64
	price           float64
	transactionType string
}

type PriorityQueue []*order

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	if pq[i].transactionType != pq[j].transactionType {
		panic("Different orderTypes in orderqueue")
	}
	if pq[i].transactionType == "BUY" {
		return pq[i].price > pq[j].price
	} else {
		return pq[i].price < pq[j].price
	}
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x any) {
	*pq = append(*pq, x.(*order))
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

func (pq PriorityQueue) Top() *order {
	return pq[len(pq)-1]
}

func (pq *PriorityQueue) List() {
	for pq.Len() > 0 {
		item := heap.Pop(pq).(*order)
		fmt.Printf("%.2f:%d ", item.price, item.quantity)
	}
}
