package redis

import redigo "github.com/garyburd/redigo/redis"

// Increment function
func (r Redis) Increment(key string) (int, error) {
	conn := r.Pool.Get()
	defer conn.Close()
	return redigo.Int(conn.Do("INCR", key))
}

// IncrementBy function
func (r Redis) IncrementBy(key string, value int) (int, error) {
	conn := r.Pool.Get()
	defer conn.Close()
	return redigo.Int(conn.Do("INCRBY", key, value))
}

// Dcrement function
func (r Redis) Dcrement(key string) (int, error) {
	conn := r.Pool.Get()
	defer conn.Close()
	return redigo.Int(conn.Do("DECR", key))
}

// DcrementBy function
func (r Redis) DcrementBy(key string, value int) (int, error) {
	conn := r.Pool.Get()
	defer conn.Close()
	return redigo.Int(conn.Do("DECRBY", key, value))
}
