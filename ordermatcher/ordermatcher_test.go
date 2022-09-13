package ordermatcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// test function
func TestOrderMatcher(t *testing.T) {
	om := OrderMatcher{}
	om.Start()

	om.AddBuy("niffy", 100, 4.0)
	om.AddBuy("niffy", 200, 4.0)
	om.AddBuy("niffy", 300, 4.0)
	assert.Equal(t, om.orderQueues["niffy:BUY"].Len(), 3)
}
