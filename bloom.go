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

func New(m uint, k uint, bitSet BitSetProvider) *BloomFilter {
	return &BloomFilter{m: m, k: k, bitSet: bitSet}
}

func EstimateParameters(n uint, p float64) (uint, uint) {
	m := math.Ceil(float64(n) * math.Log(p) / math.Log(1.0/math.Pow(2.0, math.Ln2)))
	k := math.Ln2*m/float64(n) + 0.5

	return uint(m), uint(k)
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
	locations := make([]uint, f.k)
	hasher := fnv.New64()
	hasher.Write(data)
	a := make([]byte, 1)
	for i := uint(0); i < f.k; i++ {
		a[0] = byte(i)
		hasher.Write(a)
		hashValue := hasher.Sum64()
		locations[i] = uint(hashValue % uint64(f.m))
	}
	return locations
}
