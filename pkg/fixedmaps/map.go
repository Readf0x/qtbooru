package fixedmaps

import "sync"

type FixedSizeMap[K comparable, V any] struct {
	maxSize int
	keys    []K
	data    map[K]V
	mu      sync.Mutex
}

func NewFixedSizeMap[K comparable, V any](maxSize int) *FixedSizeMap[K, V] {
	return &FixedSizeMap[K, V]{
		maxSize: maxSize,
		keys:    make([]K, 0, maxSize),
		data:    make(map[K]V),
	}
}

func (m *FixedSizeMap[K, V]) Set(key K, value V) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.data[key]; exists {
		m.data[key] = value
		for i, k := range m.keys {
			if k == key {
				m.keys = append(m.keys[:i], m.keys[i+1:]...)
				m.keys = append(m.keys, key)
				break
			}
		}
		return
	}

	if len(m.keys) >= m.maxSize {
		oldestKey := m.keys[0]
		delete(m.data, oldestKey)
		m.keys = m.keys[1:]
	}

	m.data[key] = value
	m.keys = append(m.keys, key)
}

func (m *FixedSizeMap[K, V]) Get(key K) (V, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	value, exists := m.data[key]
	if exists {
		for i, k := range m.keys {
			if k == key {
				m.keys = append(m.keys[:i], m.keys[i+1:]...)
				m.keys = append(m.keys, key)
				break
			}
		}
	}
	return value, exists
}

func (m *FixedSizeMap[K, V]) Size() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.data)
}

func (m *FixedSizeMap[K, V]) Keys() []K {
	m.mu.Lock()
	defer m.mu.Unlock()
	return append([]K(nil), m.keys...)
}

