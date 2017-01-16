package bloom

import (
	"github.com/garyburd/redigo/redis"
)

type Connection interface {
	Do(cmd string, args ...interface{}) (reply interface{}, err error)
	Flush() error
	Send(cmd string, args ...interface{}) error
}

type RedisBitSet struct {
	key  string
	conn Connection
}

func NewRedisBitSet(key string, conn Connection) *RedisBitSet {
	return &RedisBitSet{key, conn}
}

func (r *RedisBitSet) Set(offsets []uint) error {
	for _, offset := range offsets {
		err := r.conn.Send("SETBIT", r.key, offset, 1)
		if err != nil {
			return err
		}
	}

	return r.conn.Flush()
}

func (r *RedisBitSet) Test(offsets []uint) (bool, error) {
	for _, offset := range offsets {
		bitValue, err := redis.Int(r.conn.Do("GETBIT", r.key, offset))
		if err != nil {
			return false, err
		}
		if bitValue == 0 {
			return false, nil
		}
	}

	return true, nil
}
