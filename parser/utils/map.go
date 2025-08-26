package utils

// KeyValue pair for ordered iteration
type KeyValue[K comparable, V any] struct {
	Key   K
	Value V
}

type OrderedMap[K comparable, V any] struct {
	data map[K]V
	keys []K
}

func NewOrderedMap[K comparable, V any]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{data: make(map[K]V)}
}

// Insert or update key
func (om *OrderedMap[K, V]) Set(key K, value V) {
	if _, exists := om.data[key]; !exists {
		om.keys = append(om.keys, key)
	}
	om.data[key] = value
}

// Get value
func (om *OrderedMap[K, V]) Get(key K) (V, bool) {
	val, ok := om.data[key]
	return val, ok
}

// Iter returns ordered key-value pairs
func (om *OrderedMap[K, V]) Iter() []KeyValue[K, V] {
	result := make([]KeyValue[K, V], 0, len(om.keys))
	for _, k := range om.keys {
		result = append(result, KeyValue[K, V]{Key: k, Value: om.data[k]})
	}
	return result
}

func (om *OrderedMap[K, V]) Range(f func(key K, value V) bool) {
	for _, k := range om.keys {
		v := om.data[k]
		if !f(k, v) {
			break
		}
	}
}
