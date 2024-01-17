package local

type Cache struct {
	storage map[int64]*Value
}

func New() *Cache {
	return &Cache{
		storage: map[int64]*Value{},
	}
}

func (c *Cache) SaveString(key int64, data string) bool {
	if v, ok := c.storage[key]; ok && v != nil {
		v.String = &data
	} else {
		c.storage[key] = &Value{
			String: &data,
		}
	}

	return true
}

func (c *Cache) SaveInt(key int64, data int) bool {
	if v, ok := c.storage[key]; ok && v != nil {
		v.Int = &data
	} else {
		c.storage[key] = &Value{
			Int: &data,
		}
	}

	return true
}

func (c *Cache) SaveArray(key int64, index int, data string, len *int) bool {
	if v, ok := c.storage[key]; ok {
		if v != nil && v.Array == nil && len != nil {
			v.Array = make([]string, *len)
		}
		if v == nil && len != nil {
			c.storage[key] = &Value{
				Array: make([]string, *len),
			}
		}
		c.storage[key].Array[index] = data
	} else {
		if len != nil {
			c.storage[key] = &Value{
				Array: make([]string, *len),
			}
		}
		c.storage[key].Array[index] = data
	}

	return true
}

func (c *Cache) Get(key int64) *Value {
	return c.storage[key]
}

func (c *Cache) Clear(key int64) bool {
	c.storage[key] = nil

	return true
}
