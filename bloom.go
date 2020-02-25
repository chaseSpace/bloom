package bloom

import (
	"crypto/sha256"
	"github.com/spaolacci/murmur3"
	"log"
)

type FalseJudgeRateConfig struct {
	mnRatio uint // value of m/n, m is bitmap size, n is number of elements
	k       uint // times of hash
}

type FalseRateTyp float32

// listing some specific FalseRateTyp value, they have corresponding k value
// and m/n value
const (
	OneDiv10thousand    FalseRateTyp = 1e-4
	EightDiv100thousand              = 8.53e-5
	FiveDiv1million                  = 5.73e-6
)

var _configMap = map[FalseRateTyp]*FalseJudgeRateConfig{
	OneDiv10thousand:    {28, 5},
	EightDiv100thousand: {30, 5},
	FiveDiv1million:     {32, 8},
}

// GetFalseJudgeRateConfig return config of false judge rate by the specific falseJudgeRate
func GetFalseJudgeRateConfig(falseJudgeRate FalseRateTyp) *FalseJudgeRateConfig {
	return _configMap[falseJudgeRate]
}

type Bloom interface {
	Put(elem []byte) error
	Exist(elem []byte) (bool, error)
	Close() error
	AppliedSpaceWithKBytes() float32
}

func hashFunc(data []byte, seed int) (uint64, error) {
	// two hash operations to reduce collision rate
	sha := sha256.Sum256(data)
	data = sha[:]
	m := murmur3.New128WithSeed(uint32(seed))
	_, err := m.Write(data)
	if err != nil {
		return 0, err
	}
	r, _ := m.Sum128()
	return r, nil
}

// MemBloomFilter is a representation of memory bloom filter
type MemBloomFilter struct {
	conf   *FalseJudgeRateConfig
	bitmap *BitMap
}

// NewMemBloomFilter create and return a Bloom Filter using memory
func NewMemBloomFilter(falseRateConf *FalseJudgeRateConfig, DupeElementNumber uint64) Bloom {
	bitmapCapacity := DupeElementNumber * uint64(falseRateConf.mnRatio)
	log.Printf("DupeElementNumber:%d cap:%d\n", DupeElementNumber, bitmapCapacity)
	bitmap := NewBitMap(bitmapCapacity)
	return &MemBloomFilter{falseRateConf, &bitmap}
}

// Put put a []byte element into bitmap
func (bf *MemBloomFilter) Put(elem []byte) error {
	if elem == nil {
		return nil
	}
	for i := uint(0); i < bf.conf.k; i++ {
		hashValue, err := hashFunc(elem, int(i))
		if err != nil {
			return err
		}
		bf.bitmap.Set(hashValue % bf.bitmap.Capacity())
	}
	return nil
}

// Exist returns whether the given element exists in the bitmap
// if elem is nil, it always return true
func (bf *MemBloomFilter) Exist(elem []byte) (bool, error) {
	if elem == nil {
		return true, nil
	}
	for i := uint(0); i < bf.conf.k; i++ {
		hashValue, err := hashFunc(elem, int(i))
		if err != nil {
			return false, err
		}
		if !bf.bitmap.IsSet(hashValue % bf.bitmap.Capacity()) {
			return false, nil
		}
	}
	return true, nil
}

// Close release the applied resource
func (bf *MemBloomFilter) Close() error {
	bf.bitmap.Close()
	return nil
}

// AppliedSpaceWithKBytes return the applied space size with KBytes of
// underlying bitmap
func (bf *MemBloomFilter) AppliedSpaceWithKBytes() float32 {
	return bf.bitmap.AppliedSpaceWithKBytes()
}
