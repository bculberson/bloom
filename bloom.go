package bloom

import (
	"hash/fnv"
	"math"
)

type BitSetProvider interface {
	Set(uint) error
	Test(uint) (bool, error)
}

type BloomFilter struct {
	m      uint
	k      uint
	bitSet BitSetProvider
}

// Creates a new Bloom filter for about n items with fp false positive rate
func New(n uint, fp float64, bitSet BitSetProvider) (*BloomFilter, error) {
	m, k := estimateParameters(n, fp)
	return &BloomFilter{m: m, k: k, bitSet: bitSet}, nil
}

// Used with permission from https://bitbucket.org/ww/bloom/src/829aa19d01d9/bloom.go
func estimateParameters(n uint, p float64) (uint, uint) {
	m := uint(math.Ceil(-1 * float64(n) * math.Log(p) / math.Pow(math.Log(2), 2)))
	k := uint(math.Ceil(math.Log(2) * float64(m) / float64(n)))
	return m, k
}

func (f *BloomFilter) Add(data []byte) error {
	hashValue := getHash(data)
	for i := uint(0); i < f.k; i++ {
		location := f.getLocation(hashValue, i)
		err := f.bitSet.Set(location)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *BloomFilter) Exists(data []byte) (bool, error) {
	hashValue := getHash(data)
	for i := uint(0); i < f.k; i++ {
		location := f.getLocation(hashValue, i)
		isSet, err := f.bitSet.Test(location)
		if err != nil {
			return false, err
		}
		if !isSet {
			return false, nil
		}
	}
	return true, nil
}

func (f *BloomFilter) getLocation(hashValue uint64, i uint) uint {
	return uint((hashValue * uint64(i+1)) % uint64(f.m))
}

func getHash(data []byte) uint64 {
	hasher := fnv.New64()
	hasher.Write(data)
	return hasher.Sum64()
}
