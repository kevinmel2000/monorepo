package redis

import (
	"errors"

	redigo "github.com/garyburd/redigo/redis"
)

// Get string value
func (r Redis) Get(key string) (string, error) {
	conn := r.Pool.Get()
	defer conn.Close()
	v, err := redigo.String(conn.Do("GET", key))
	return v, err
}

// Set key and value
func (r Redis) Set(key, value string) (string, error) {
	conn := r.Pool.Get()
	defer conn.Close()
	return redigo.String(conn.Do("SET", key, value))
}

// SetWithNX with NX params
func (r Redis) SetWithNX(key, value string, expire int) (string, error) {
	conn := r.Pool.Get()
	defer conn.Close()
	return redigo.String(conn.Do("SET", key, value, "NX", "EX", expire))
}

// SetNX key and value
func (r Redis) SetNX(key, value string, expire int) (int, error) {
	conn := r.Pool.Get()
	defer conn.Close()
	return redigo.Int(conn.Do("SETNX", key, value))
}

// SetEX key and value
func (r Redis) SetEX(key, value string, expire int) (string, error) {
	conn := r.Pool.Get()
	defer conn.Close()
	return redigo.String(conn.Do("SETEX", key, expire, value))
}

// MGet keys and value
func (r Redis) MGet(keys []string) ([]string, error) {
	conn := r.Pool.Get()
	defer conn.Close()

	var keysItfs []interface{}
	for _, key := range keys {
		keysItfs = append(keysItfs, key)
	}

	return redigo.Strings(conn.Do("MGET", keysItfs...))
}

// HSet key and value
func (r Redis) HSet(key, child, value string, expire int) (string, error) {
	conn := r.Pool.Get()
	defer conn.Close()
	_, err := redigo.String(conn.Do("HSET", key, child, value))
	if err != nil {
		return "", err
	}
	return redigo.String(conn.Do("EXPIRE", key, expire))
}

// HGet key and value
func (r Redis) HGet(key, child string) (string, error) {
	conn := r.Pool.Get()
	defer conn.Close()
	return redigo.String(conn.Do("HGET", key, child))
}

// HMSet function
func (r Redis) HMSet(key string, childs []string, vals []string, expire int) (string, error) {
	conn := r.Pool.Get()
	defer conn.Close()

	var valueItfs []interface{}
	if len(childs) != len(vals) {
		return "", errors.New("Len isnt same")
	}

	valueItfs = append(valueItfs, key)
	for i, child := range childs {
		valueItfs = append(valueItfs, child)
		valueItfs = append(valueItfs, vals[i])
	}
	return redigo.String(conn.Do("HMSET", valueItfs...))
}

// HMGet keys and value
func (r Redis) HMGet(key string, childs []string) ([]string, error) {
	conn := r.Pool.Get()
	defer conn.Close()

	var keyItfs []interface{}
	keyItfs = append(keyItfs, key)
	for _, child := range childs {
		keyItfs = append(keyItfs, child)
	}
	return redigo.Strings(conn.Do("HMGET", keyItfs...))
}
