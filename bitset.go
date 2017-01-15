package bloom

import (
	"github.com/willf/bitset"
)

type BitSet struct {
	bitSet *bitset.BitSet
}

func NewBitSet() *BitSet {
	return &BitSet{}
}

func (b *BitSet) Set(offsets []uint) error {
	maxOffset := uint(0)
	for _, offset := range offsets {
		if offset > maxOffset {
			maxOffset = offset
		}
	}
	if b.bitSet == nil {
		b.bitSet = bitset.New(maxOffset)
	}
	if maxOffset > b.bitSet.Len() {
		previousBitSet := b.bitSet
		b.bitSet = bitset.New(maxOffset)
		b.bitSet.InPlaceUnion(previousBitSet)
	}
	for _, offset := range offsets {
		b.bitSet.Set(offset)
	}
	return nil
}

func (b *BitSet) Test(offsets []uint) (bool, error) {
	if b.bitSet == nil {
		return false, nil
	}
	for _, offset := range offsets {
		if offset > b.bitSet.Len() || !b.bitSet.Test(offset) {
			return false, nil
		}
	}

	return true, nil
}
