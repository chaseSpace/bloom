package bloom

import (
	"testing"
)

func TestBitMap(t *testing.T) {
	var capacity uint64 = 1 << 20
	var bm = NewBitMap(uint64(capacity))

	var oddNumberPositions []int
	var evenNumberPositions []int
	for i := 0; uint64(i) < capacity; i++ {
		if i%2 == 0 {
			bm.Set(uint64(i))
			evenNumberPositions = append(evenNumberPositions, i)
		} else {
			oddNumberPositions = append(oddNumberPositions, i)
		}
	}

	for _, i := range evenNumberPositions {
		if !bm.IsSet(uint64(i)) {
			t.Fatalf("even number position should be set, pos:%d", i)
		}
	}

	for _, i := range oddNumberPositions {
		if bm.IsSet(uint64(i)) {
			t.Fatalf("odd number position shouldn't be set, pos:%d", i)
		}
	}
}
