package bloom

import (
	"errors"

	"github.com/willf/bitset"
)

var ErrBitSetUninitialized = errors.New("BitSet Uninitialized")

type BitSet struct {
	bitSet *bitset.BitSet
	size   uint
}

func NewBitSet() *BitSet {
	return &BitSet{}
}

func (b *BitSet) New(size uint) error {
	b.bitSet = bitset.New(size)
	b.size = size
	return nil
}

func (b *BitSet) Set(offset uint) error {
	if b.bitSet == nil {
		return ErrBitSetUninitialized
	}
	b.bitSet.Set(offset)
	return nil
}

func (b *BitSet) Test(offset uint) (bool, error) {
	if b.bitSet == nil {
		return false, ErrBitSetUninitialized
	}
	return b.bitSet.Test(offset), nil
}
