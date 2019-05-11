package dw

import (
	"sort"
)

type Int32Array []int32

func (arr Int32Array) Len() int {
	return len(arr)
}

func (arr Int32Array) Less(i, j int) bool {
	return arr[i] < arr[j]
}

func (arr Int32Array) Swap(i, j int) {
	arr[i], arr[j] = arr[j], arr[i]
}

func SortInt32Array(arr Int32Array) {
	sort.Sort(arr)
}
