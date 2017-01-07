# Bloom Filters for Golang

Bloom filter for go, backed by redis or in process bitset

For more information on the uses of Bloom filters, please read
https://en.wikipedia.org/wiki/Bloom_filter

## Example Usage (in process):

This bloom filter is initialized to hold 1000 keys and
will have a false positive rate of 1% (.01).

```go
b, _ := bloom.New(1000, .01, bloom.NewBitSet())
b.Add([]byte("some key"))
exists, _ := b.Exists([]byte("some key"))
doesNotExist, _ := b.Exists([]byte("some other key"))
```

## Example Usage (redis backed):

This bloom filter is initialized to hold 1000 keys and
will have a false positive rate of 1% (.01).

This library uses http://github.com/garyburd/redigo/redis

```go
pool := &redis.Pool{
    MaxIdle:     3,
    IdleTimeout: 240 * time.Second,
    Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", addr) },
}
	
bitSet := bloom.NewRedisBitSet("test_key", pool)
b, _ := bloom.New(1000, .01, bitSet)
b.Add([]byte("some key"))
exists, _ := b.Exists([]byte("some key"))
doesNotExist, _ := b.Exists([]byte("some other key"))
```

