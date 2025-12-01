package sortedmap

// SortedMapByValue: K = map key, V = map value
// 排序依据为 cmpValue(v1, v2)
type SortedMapByValue[K comparable, V any] struct {
	m        map[K]V
	keys     []K
	cmpValue func(a, b V) int // value comparator
}

func NewSortedMapByValue[K comparable, V any](cmpValue func(a, b V) int) *SortedMapByValue[K, V] {
	return &SortedMapByValue[K, V]{
		m:        make(map[K]V),
		keys:     make([]K, 0),
		cmpValue: cmpValue,
	}
}

// 二分查找：按 value 排序查找插入位置
func (s *SortedMapByValue[K, V]) binarySearch(v V) int {
	lo, hi := 0, len(s.keys)
	for lo < hi {
		mid := (lo + hi) / 2
		kv := s.m[s.keys[mid]]
		if s.cmpValue(kv, v) < 0 { // kv < v
			lo = mid + 1
		} else {
			hi = mid
		}
	}
	return lo
}

func (s *SortedMapByValue[K, V]) Put(k K, v V) {
	_, exists := s.m[k]
	if exists {
		// key 已存在，要移除旧位置
		s.Delete(k)
	}

	// 在 value 排序序列中，找到插入位置
	pos := s.binarySearch(v)

	// 插入 keys[pos]
	s.keys = append(s.keys, k)
	copy(s.keys[pos+1:], s.keys[pos:])
	s.keys[pos] = k

	// 更新 map
	s.m[k] = v
}

func (s *SortedMapByValue[K, V]) Get(k K) (V, bool) {
	v, ok := s.m[k]
	return v, ok
}

func (s *SortedMapByValue[K, V]) Keys() []K {
	out := make([]K, len(s.keys))
	copy(out, s.keys)
	return out
}

func (s *SortedMapByValue[K, V]) Values() []V {
	out := make([]V, len(s.keys))
	for i, k := range s.keys {
		out[i] = s.m[k]
	}
	return out
}

func (s *SortedMapByValue[K, V]) Range(fn func(k K, v V) bool) {
	for _, k := range s.keys {
		if !fn(k, s.m[k]) {
			return
		}
	}
}

func (s *SortedMapByValue[K, V]) Delete(k K) {
	_, exists := s.m[k]
	if !exists {
		return
	}
	delete(s.m, k)

	// 按 key 线性扫描删除
	for i := range s.keys {
		if s.keys[i] == k {
			copy(s.keys[i:], s.keys[i+1:])
			s.keys = s.keys[:len(s.keys)-1]
			return
		}
	}
}
