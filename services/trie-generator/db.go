package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/garyburd/redigo/redis"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
)

var (
	Pool *redis.Pool
)

func newPool(server string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,

		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			return c, err
		},

		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

type RedisDB struct {
	dbPool *redis.Pool
	lock   sync.RWMutex
}

func NewRedisDB(server string) *RedisDB {
	r := &RedisDB{
		dbPool: newPool(server),
	}

	return r
}

func (r *RedisDB) Keys() [][]byte {

	// log.Printf("Keys()\n")

	conn := r.dbPool.Get()
	defer conn.Close()

	keys, err := redis.ByteSlices(conn.Do("KEYS", "*"))
	if err != nil {
		log.Printf("Error getting value in redisDB: %v", err)
		return nil
	}

	response := [][]byte{}
	for _, key := range keys {
		response = append(response, []byte(key))
	}

	return response

}

func (r *RedisDB) Put(key []byte, value []byte) error {

	// log.Printf("Put %x %x\n", key, value)

	conn := r.dbPool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, value)
	if err != nil {
		return fmt.Errorf("Error setting value in redisDB: %v", err)
	}

	return nil
}

func (r *RedisDB) Get(key []byte) ([]byte, error) {

	// log.Printf("Get %x\n", key)

	conn := r.dbPool.Get()
	defer conn.Close()

	value, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		fmt.Printf("Error getting value from redisDB: %v\n", err)
		return nil, err
	}

	return value, nil
}

func (r *RedisDB) Has(key []byte) (bool, error) {

	// log.Printf("Has %x\n", key)

	conn := r.dbPool.Get()
	defer conn.Close()

	_, err := redis.ByteSlices(conn.Do("GET", key))
	if err != nil {
		if err != redis.ErrNil {
			fmt.Printf("Error getting value from redisDB: %v\n", err)
			return false, err
		} else {
			// value not found
			return false, nil
		}
	}

	return true, nil
}

func (r *RedisDB) Delete(key []byte) error {

	// log.Printf("Delete() %x\n", key)

	conn := r.dbPool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", key)
	if err != nil {
		return fmt.Errorf("Error deleting value in redisDB: %v", err)
	}

	return nil
}

func (r *RedisDB) Close() {

	// are we implementing this one in here?

	// log.Printf("Close()\n")

}

func (r *RedisDB) NewBatch() ethdb.Batch {

	// log.Printf("NewBatch()\n")

	return &redisBatch{db: r}
}

type kv struct{ k, v []byte }

type redisBatch struct {
	db     *RedisDB
	writes []kv
	size   int
}

func (r *redisBatch) Put(key, value []byte) error {

	// log.Printf("memBatch.Put() %x %x\n", key, value)

	r.writes = append(r.writes, kv{common.CopyBytes(key), common.CopyBytes(value)})
	r.size += len(value)
	return nil
}

func (r *redisBatch) Write() error {

	// log.Printf("memBatch.Write()\n")

	r.db.lock.Lock()
	defer r.db.lock.Unlock()

	for _, kv := range r.writes {
		r.db.Put(kv.k, kv.v)
	}
	return nil
}

func (r *redisBatch) ValueSize() int {

	// log.Printf("memBatch.ValueSize()")

	return r.size
}

func (r *redisBatch) Reset() {

	// log.Printf("memBatch.Reset()\n")

	r.writes = r.writes[:0]
	r.size = 0
}
