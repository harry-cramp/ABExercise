package store

import "sync/atomic"

var Quantity atomic.Int64
var TicketNumber atomic.Int64

func GetQuantity() int64 {
    return Quantity.Load()
}

func AttemptBuy() bool {
    for {
        old := Quantity.Load()
        if old == 0 {
            return false
        }
        if Quantity.CompareAndSwap(old, old-1) {
            return true
        }
    }
}

func GetTicketNumber() int64 {
	return TicketNumber.Load()
}
