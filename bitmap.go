package bloom

type BitMap struct {
	int64slice []int64
	capacity   uint64
}

// NewBitMap return a new BitMap with a specific capacity
func NewBitMap(cap uint64) BitMap {
	var bm = BitMap{make([]int64, cap/64+1), (cap/64 + 1) * 64}
	return bm
}

// Set set the bitmap corresponding position to 1
func (b BitMap) Set(pos uint64) {
	b.int64slice[(pos / 64)] |= 1 << (pos % 64)
}

// IsSet return whether the bitmap corresponding position is set to 1
func (b BitMap) IsSet(pos uint64) bool {
	var bitMapNumber = b.int64slice[(pos / 64)]
	return bitMapNumber|(1<<(pos%64)) == bitMapNumber
}

// UnSet set the bitmap corresponding position to 0
// Note: bloom filter will not use this method
func (b BitMap) UnSet(pos uint64) {
	if b.IsSet(pos) {
		b.int64slice[(pos / 64)] ^= 1 << (pos % 64)
	}
}

func (b BitMap) Capacity() uint64 {
	return b.capacity
}

func (b BitMap) AppliedSpaceWithKBytes() float32 {
	return float32(b.Capacity()) / 8 / 1024
}

func (b BitMap) Close() {
	b.int64slice = nil
}
