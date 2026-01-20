package kvtypes

type InvokeKVOptions struct {
	EX int  `json:"EX"`
	NX bool `json:"NX"`
	XX bool `json:"XX"`
}

type InvokeKVRequest struct {
	RID     string          `json:"rid"`
	Key     string          `json:"key"`
	Value   string          `json:"value"`
	Method  string          `json:"method"`
	Options InvokeKVOptions `json:"options"`
	Offset  int             `json:"offset"`
	Size    int             `json:"size"`
}

type IKVStorage interface {
	Put(bucket string, key string, value []byte, ttl int) (int, error)
	PutNX(bucket string, key string, value []byte, ttl int) (int, error)
	PutXX(bucket string, key string, value []byte, ttl int) (int, error)
	Get(bucket string, key string) ([]byte, error)
	Del(bucket string, key string) error
	Keys(bucket string, prefix string, offset int, size int) ([]string, error)
	Close()
}
