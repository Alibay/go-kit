package go_kit

func ConvertSlice[TSrc any, TRes any](src []*TSrc, converter func(*TSrc) *TRes) []*TRes {
	r := make([]*TRes, 0, len(src))
	for _, i := range src {
		r = append(r, converter(i))
	}
	return r
}

func GroupBy[TItem any, TKey comparable](slice []TItem, keyFn func(TItem) TKey) map[TKey][]TItem {
	r := make(map[TKey][]TItem)
	for _, i := range slice {
		r[keyFn(i)] = append(r[keyFn(i)], i)
	}
	return r
}

func Map[TItem any, TRes any](slice []TItem, mapFn func(TItem) TRes) []TRes {
	r := make([]TRes, 0, len(slice))
	for _, i := range slice {
		r = append(r, mapFn(i))
	}
	return r
}

func Filter[TItem any](slice []TItem, filterFn func(TItem) bool) []TItem {
	r := make([]TItem, 0, len(slice))
	for _, i := range slice {
		if filterFn(i) {
			r = append(r, i)
		}
	}
	return r
}

func Reduce[TItem any, TRes any, TKey comparable](slice []TItem, grpFn func(TItem) TKey, accFn func(TItem, TRes) TRes) map[TKey]TRes {
	r := make(map[TKey]TRes)
	for _, i := range slice {
		r[grpFn(i)] = accFn(i, r[grpFn(i)])
	}
	return r
}

func ForAll[TItem any](slice []TItem, fn func(TItem)) []TItem {
	for _, i := range slice {
		fn(i)
	}
	return slice
}

func First[TItem any](slice []*TItem, selectFn func(*TItem) bool) *TItem {
	for _, i := range slice {
		if selectFn(i) {
			return i
		}
	}
	return GetDefault[*TItem]()
}

func GetDefault[T any]() T {
	var result T
	return result
}
