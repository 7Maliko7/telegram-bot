package redis

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/redis/go-redis/v9"

	"github.com/7Maliko7/telegram-bot/pkg/cache"
)

type Cache struct {
	storage *redis.Client
}

func New(dsn string) (*Cache, error) {
	opt, err := redis.ParseURL(dsn)
	if err != nil {
		return nil, err
	}

	return &Cache{
		storage: redis.NewClient(opt),
	}, nil
}

func (c *Cache) Close() error {
	return c.storage.Close()
}

func (c *Cache) SaveInt(key int64, data int) (bool, error) {
	var (
		err error
		val []byte
	)
	if c.storage.Exists(context.TODO(), strconv.FormatInt(key, 10)).Val() == 0 {
		val, err = json.Marshal(cache.Value{
			Int: &data,
		})
		if err != nil {
			return false, err
		}
	} else {
		oldValByte, err := c.storage.Get(context.TODO(), strconv.FormatInt(key, 10)).Bytes()
		if err != nil {
			return false, err
		}
		tempVal := cache.Value{}
		err = json.Unmarshal(oldValByte, &tempVal)
		if err != nil {
			return false, err
		}

		tempVal.Int = &data
		val, err = json.Marshal(tempVal)
		if err != nil {
			return false, err
		}
	}

	err = c.storage.Set(context.TODO(), strconv.FormatInt(key, 10), string(val), 0).Err()
	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *Cache) SaveString(key int64, data string) (bool, error) {
	var (
		err error
		val []byte
	)
	if c.storage.Exists(context.TODO(), strconv.FormatInt(key, 10)).Val() == 0 {
		val, err = json.Marshal(cache.Value{
			String: &data,
		})
		if err != nil {
			return false, err
		}
	} else {
		oldValByte, err := c.storage.Get(context.TODO(), strconv.FormatInt(key, 10)).Bytes()
		if err != nil {
			return false, err
		}
		tempVal := cache.Value{}
		err = json.Unmarshal(oldValByte, &tempVal)
		if err != nil {
			return false, err
		}

		tempVal.String = &data
		val, err = json.Marshal(tempVal)
		if err != nil {
			return false, err
		}
	}

	err = c.storage.Set(context.TODO(), strconv.FormatInt(key, 10), string(val), 0).Err()
	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *Cache) SaveArray(key int64, index int, data string, len *int) (bool, error) {
	var (
		err error
		val []byte
	)
	if c.storage.Exists(context.TODO(), strconv.FormatInt(key, 10)).Val() == 0 {
		arr := cache.Value{
			Array: make([]string, *len),
		}
		arr.Array[index] = data
		val, err = json.Marshal(arr)
		if err != nil {
			return false, err
		}
	} else {
		oldValByte, err := c.storage.Get(context.TODO(), strconv.FormatInt(key, 10)).Bytes()
		if err != nil {
			return false, err
		}
		tempVal := cache.Value{}
		err = json.Unmarshal(oldValByte, &tempVal)
		if err != nil {
			return false, err
		}

		tempVal.Array[index] = data
		val, err = json.Marshal(tempVal)
		if err != nil {
			return false, err
		}
	}

	err = c.storage.Set(context.TODO(), strconv.FormatInt(key, 10), string(val), 0).Err()
	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *Cache) Get(key int64) (*cache.Value, error) {
	if c.storage.Exists(context.TODO(), strconv.FormatInt(key, 10)).Val() == 1 {
		valByte, err := c.storage.Get(context.TODO(), strconv.FormatInt(key, 10)).Bytes()
		if err != nil {
			return nil, err
		}

		tempVal := cache.Value{}
		err = json.Unmarshal(valByte, &tempVal)
		if err != nil {
			return nil, err
		}

		return &tempVal, nil
	}

	return nil, nil
}

func (c *Cache) Clear(key int64) (bool, error) {
	if c.storage.Exists(context.TODO(), strconv.FormatInt(key, 10)).Val() == 1 {
		err := c.storage.Del(context.TODO(), strconv.FormatInt(key, 10)).Err()
		if err != nil {
			return false, err
		}

		return true, nil
	}

	return false, nil
}
