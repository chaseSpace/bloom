package bloom

import (
	"encoding/binary"
	"log"
	"testing"
)

func doTestMemBloomFilter(t *testing.T, falseRateType FalseRateTyp, dupeElemNumber uint64) {
	falseJudgeRateConf := GetFalseJudgeRateConfig(falseRateType)
	bf := NewMemBloomFilter(falseJudgeRateConf, uint64(dupeElemNumber))
	defer bf.Close()
	log.Printf("begin---[dupeElemNumber:%d max-falseRatetype:%.6f AppliedSpaceWithKBytes:%.4f] \n",
		dupeElemNumber, falseRateType, bf.AppliedSpaceWithKBytes())
	var oddNumberSlice []uint64
	var evenNumberSlice []uint64
	for i := uint64(0); i < dupeElemNumber; i++ {
		if i%2 == 0 {
			buf := make([]byte, 8)
			binary.BigEndian.PutUint64(buf, uint64(i))
			err := bf.Put(buf)
			if err != nil {
				t.Fatalf("[falseRatetype:%.6f] put err %v", falseRateType, err)
			}
			evenNumberSlice = append(evenNumberSlice, i)
		} else {
			oddNumberSlice = append(oddNumberSlice, i)
		}
	}

	testFalseNumber := 0
	for i := uint64(0); i < dupeElemNumber; i++ {
		buf := make([]byte, 8)
		binary.BigEndian.PutUint64(buf, uint64(i))
		has, err := bf.Exist(buf)
		if err != nil {
			t.Fatalf("exist err %v", err)
		}
		if i%2 == 0 {
			if !has {
				t.Fatalf("[falseRatetype:%.6f] evenNumberSlice: %d not exist?\n", falseRateType, i)
			}
		} else {
			if has {
				testFalseNumber++
			}
		}
	}

	testFalseJudgeRate := float32(testFalseNumber) / float32(dupeElemNumber)

	if testFalseJudgeRate > float32(falseRateType) {
		t.Fatalf("o(╥﹏╥)o, bigger falseJudgeRate %.6f > %.6f",
			testFalseJudgeRate, falseRateType)
	} else {
		log.Printf("(>‿◠)✌, test passed! testFalseJudgeRate:%.6f  expected:%.6f\n",
			testFalseJudgeRate, falseRateType)
	}
}

func TestMemBloomFilter(t *testing.T) {
	// with TDT test
	type tableRow struct {
		dupeElemNumber uint64
		FalseRateTyp
	}
	// 1 << 20 = 1048576
	tableTestDatas := []tableRow{
		{1 << 20, OneDiv10thousand},
		{1 << 20, EightDiv100thousand},
		{1 << 20, FiveDiv1million},
	}
	//var wg sync.WaitGroup
	for _, row := range tableTestDatas {
		doTestMemBloomFilter(t, row.FalseRateTyp, row.dupeElemNumber)
	}
}
