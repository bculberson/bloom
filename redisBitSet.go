package bloom

import (
	"github.com/garyburd/redigo/redis"
)

type RedisBitSet struct {
	key   string
	pool  *redis.Pool
}

func NewRedisBitSet(key string, pool *redis.Pool) (*RedisBitSet) {
	return &RedisBitSet{key, pool}
}

func (r *RedisBitSet) New(size uint) error {
	conn := r.pool.Get()
	defer conn.Close()

	return conn.Send("SETBIT", r.key, size, 0)
}

func (r *RedisBitSet) Set(offset uint) error {
	conn := r.pool.Get()
	defer conn.Close()

	return conn.Send("SETBIT", r.key, offset, 1)
}

func (r *RedisBitSet) Test(offset uint) (bool, error) {
	conn := r.pool.Get()
	defer conn.Close()

	bitValue, err := redis.Int(conn.Do("GETBIT", r.key, offset))
	if err != nil {
		return false, err
	}

	return bitValue == 1, nil
}
