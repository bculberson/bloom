package bloom_test

import (
	"testing"
	"time"

	"encoding/binary"
	"github.com/alicebob/miniredis"
	"github.com/bculberson/bloom"
	"github.com/garyburd/redigo/redis"
)

func TestRedisBloomFilter(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Error("Miniredis could not start")
	}
	defer s.Close()

	pool := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", s.Addr()) },
	}
	conn := pool.Get()
	defer conn.Close()

	bitSet := bloom.NewRedisBitSet("test_key", conn)
	b := bloom.New(1000, .01, bitSet)
	testBloomFilter(t, b)
}

func TestBloomFilter(t *testing.T) {
	b := bloom.New(1000, .01, bloom.NewBitSet())
	testBloomFilter(t, b)
}

func TestCollision(t *testing.T) {
	b := bloom.New(100, .01, bloom.NewBitSet())
	shouldNotExist := 0
	for i := 0; i < 100; i++ {
		data := make([]byte, 4)
		binary.LittleEndian.PutUint32(data, uint32(i))
		existsBefore, err := b.Exists(data)
		if err != nil {
			t.Fatal("Error checking existence.")
		}
		if existsBefore {
			shouldNotExist = shouldNotExist + 1
		}
		err = b.Add(data)
		if err != nil {
			t.Fatal("Error adding item.")
		}
		existsAfter, err := b.Exists(data)
		if err != nil {
			t.Fatal("Error checking existence.")
		}
		if !existsAfter {
			t.Fatal("Item should exist.")
		}
	}
	if shouldNotExist > 2 {
		t.Fatal("Too many false positives.")
	}
}

func testBloomFilter(t *testing.T, b *bloom.BloomFilter) {
	data := []byte("some key")
	existsBefore, err := b.Exists(data)
	if err != nil {
		t.Fatal("Error checking for existence in bloom filter")
	}
	if existsBefore {
		t.Fatal("Bloom filter should not contain this data")
	}
	err = b.Add(data)
	if err != nil {
		t.Fatal("Error adding to bloom filter")
	}
	existsAfter, err := b.Exists(data)
	if err != nil {
		t.Fatal("Error checking for existence in bloom filter")
	}
	if !existsAfter {
		t.Fatal("Bloom filter should contain this data")
	}
}
