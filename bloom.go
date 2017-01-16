package bloom

import (
	"hash/fnv"
	"math"
)

type BitSetProvider interface {
	Set([]uint) error
	Test([]uint) (bool, error)
}

type BloomFilter struct {
	m      uint
	k      uint
	bitSet BitSetProvider
}

// Creates a new Bloom filter for about n items with fp false positive rate
func New(n uint, fp float64, bitSet BitSetProvider) *BloomFilter {
	m, k := estimateParameters(n, fp)
	return &BloomFilter{m: m, k: k, bitSet: bitSet}
}

// Used with permission from https://bitbucket.org/ww/bloom/src/829aa19d01d9/bloom.go
func estimateParameters(n uint, p float64) (uint, uint) {
	m := uint(math.Ceil(-1 * float64(n) * math.Log(p) / math.Pow(math.Log(2), 2)))
	k := uint(math.Ceil(math.Log(2) * float64(m) / float64(n)))
	return m, k
}

func (f *BloomFilter) Add(data []byte) error {
	locations := f.getLocations(data)
	err := f.bitSet.Set(locations)
	if err != nil {
		return err
	}
	return nil
}

func (f *BloomFilter) Exists(data []byte) (bool, error) {
	locations := f.getLocations(data)
	isSet, err := f.bitSet.Test(locations)
	if err != nil {
		return false, err
	}
	if !isSet {
		return false, nil
	}

	return true, nil
}

func (f *BloomFilter) getLocations(data []byte) []uint {
	hashValue := getHash(data)
	locations := make([]uint, f.k)
	for i := uint(0); i < f.k; i++ {
		locations[i] = uint((hashValue * uint64(i+1)) % uint64(f.m))
	}
	return locations
}

func getHash(data []byte) uint64 {
	hasher := fnv.New64()
	hasher.Write(data)
	return hasher.Sum64()
}
