package bloom

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"strings"
)

const redisMaxLength = 8 * 512 * 1024 * 1024

type Connection interface {
	Do(cmd string, args ...interface{}) (reply interface{}, err error)
	Send(cmd string, args ...interface{}) error
	Flush() error
}

type RedisBitSet struct {
	keyPrefix string
	conn      Connection
	m         uint
}

func NewRedisBitSet(keyPrefix string, m uint, conn Connection) *RedisBitSet {
	return &RedisBitSet{keyPrefix, conn, m}
}

func (r *RedisBitSet) Set(offsets []uint) error {
	for _, offset := range offsets {
		key, thisOffset := r.getKeyOffset(offset)
		err := r.conn.Send("SETBIT", key, thisOffset, 1)
		if err != nil {
			return err
		}
	}

	return r.conn.Flush()
}

func (r *RedisBitSet) Test(offsets []uint) (bool, error) {
	for _, offset := range offsets {
		key, thisOffset := r.getKeyOffset(offset)
		bitValue, err := redis.Int(r.conn.Do("GETBIT", key, thisOffset))
		if err != nil {
			return false, err
		}
		if bitValue == 0 {
			return false, nil
		}
	}

	return true, nil
}

func (r *RedisBitSet) Expire(seconds uint) error {
	n := uint(0)
	for n <= uint(r.m/redisMaxLength) {
		key := fmt.Sprintf("%s:%d", r.keyPrefix, n)
		n = n + 1
		err := r.conn.Send("EXPIRE", key, seconds)
		if err != nil {
			return err
		}
	}
	return r.conn.Flush()
}

func (r *RedisBitSet) Delete() error {
	n := uint(0)
	keys := make([]string, 0)
	for n <= uint(r.m/redisMaxLength) {
		key := fmt.Sprintf("%s:%d", r.keyPrefix, n)
		keys = append(keys, key)
		n = n + 1
	}
	_, err := r.conn.Do("DEL", strings.Join(keys, " "))
	return err
}

func (r *RedisBitSet) getKeyOffset(offset uint) (string, uint) {
	n := uint(offset / redisMaxLength)
	thisOffset := offset - n*redisMaxLength
	key := fmt.Sprintf("%s:%d", r.keyPrefix, n)
	return key, thisOffset
}
