package arr

import "fmt"

func ArrMap[T any, K any](vals []T, cb func(T) K) []K {
	result := make([]K, len(vals))
	for i, v := range vals {
		result[i] = cb(v)
	}

	return result
}

func ArrFilter[T any](vals []T, cb func(T) bool) []T {
	result := []T{}
	for _, v := range vals {
		ok := cb(v)
		if !ok {
			continue
		}

		result = append(result, v)
	}

	return result
}

func ArrEach[T any](vals []T, cb func(T)) {
	for _, v := range vals {
		cb(v)
	}
}

func ArrEachIdx[T any](vals []T, cb func(T, int)) {
	for i, v := range vals {
		cb(v, i)
	}
}

func ArrUnique[T comparable](vals []T) []T {
	valsUnique := []T{}
	valsMap := map[T]bool{}
	for _, val := range vals {
		if valsMap[val] {
			continue
		}

		valsUnique = append(valsUnique, val)
		valsMap[val] = true
	}

	return valsUnique
}

func ArrPrepend[T comparable](vals []T, v T) []T {
	if len(vals) == 0 {
		return vals
	}

	var noop T
	vals = append(vals, noop)
	copy(vals[1:], vals)
	vals[0] = v
	return vals
}

func ArrSlice[T any](items []T, from, to int) []T {
	if from < 0 {
		from = 0
	}
	if to > len(items)-1 {
		to = len(items) - 1
	}

	if from > to {
		return []T{}
	}

	return items[from : to+1]
}

func ArrAny[T any](items []T) []any {
	result := make([]any, len(items))
	for i, v := range items {
		result[i] = v
	}
	return result
}

func ArrIntersect[T comparable](arrays ...[]T) ([]T, error) {
	hash := map[T]*int{}
	for _, slice := range arrays {
		if slice == nil {
			return nil, fmt.Errorf("intersect: array item should not be nil")
		}
		slice = ArrUnique(slice)
		for _, item := range slice {
			if counter := hash[item]; counter != nil {
				*counter++
			} else {
				temp := 1
				hash[item] = &temp
			}
		}
	}

	result := make([]T, 0)
	length := len(arrays)
	for value, counter := range hash {
		if *counter == length {
			result = append(result, value)
		}
	}
	return result, nil
}

func ArrFind[T comparable](array []T, t T) (x T, found bool) {
	for _, item := range array {
		if item == t {
			return t, true
		}
	}

	return
}

func ArrCompare[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func ArrCompareRelaxed[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for _, v := range a {
		if _, found := ArrFind(b, v); !found {
			return false
		}
	}
	return true
}
