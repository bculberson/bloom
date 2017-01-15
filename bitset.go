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

func (b *BitSet) Set(offset uint) error {
	if b.bitSet == nil {
		b.bitSet = bitset.New(offset)
	}
	if offset > b.bitSet.Len() {
		previousBitSet := b.bitSet
		b.bitSet = bitset.New(offset)
		b.bitSet.InPlaceUnion(previousBitSet)
	}
	b.bitSet.Set(offset)
	return nil
}

func (b *BitSet) Test(offset uint) (bool, error) {
	if b.bitSet == nil {
		b.bitSet = bitset.New(offset)
	}
	if offset > b.bitSet.Len() {
		previousBitSet := b.bitSet
		b.bitSet = bitset.New(offset)
		b.bitSet.InPlaceUnion(previousBitSet)
	}

	return b.bitSet.Test(offset), nil
}
