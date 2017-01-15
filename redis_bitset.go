package bloom

import (
	"github.com/garyburd/redigo/redis"
)

type RedisBitSet struct {
	key  string
	pool *redis.Pool
}

func NewRedisBitSet(key string, pool *redis.Pool) *RedisBitSet {
	return &RedisBitSet{key, pool}
}

func (r *RedisBitSet) Set(offsets []uint) error {
	conn := r.pool.Get()
	defer conn.Close()

	for _, offset := range offsets {
		err := conn.Send("SETBIT", r.key, offset, 1)
		if err != nil {
			return err
		}
	}

	return conn.Flush()
}

func (r *RedisBitSet) Test(offsets []uint) (bool, error) {
	conn := r.pool.Get()
	defer conn.Close()

	for _, offset := range offsets {
		bitValue, err := redis.Int(conn.Do("GETBIT", r.key, offset))
		if err != nil {
			return false, err
		}
		if bitValue == 0 {
			return false, nil
		}
	}

	return true, nil
}
