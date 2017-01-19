package bloom

import (
	"github.com/willf/bitset"
)

type BitSet struct {
	bitSet *bitset.BitSet
}

func NewBitSet(m uint) *BitSet {
	return &BitSet{bitSet: bitset.New(m)}
}

func (b *BitSet) Set(offsets []uint) error {
	for _, offset := range offsets {
		b.bitSet.Set(offset)
	}
	return nil
}

func (b *BitSet) Test(offsets []uint) (bool, error) {
	for _, offset := range offsets {
		if !b.bitSet.Test(offset) {
			return false, nil
		}
	}

	return true, nil
}
