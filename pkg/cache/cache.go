package cache

type Cacher interface {
	Close() error
	SaveInt(key int64, data int) (bool, error)
	SaveString(key int64, data string) (bool, error)
	SaveArray(key int64, index int, data string, len *int) (bool, error)
	Get(key int64) (*Value, error)
	Clear(key int64) (bool, error)
}

type Value struct {
	Int    *int     `json:"int,omitempty"`
	String *string  `json:"string,omitempty"`
	Array  []string `json:"array,omitempty"`
}
