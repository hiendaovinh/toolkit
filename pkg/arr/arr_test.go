package arr_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/hiendaovinh/toolkit/v2/pkg/arr"
	"github.com/stretchr/testify/assert"
)

func TestArrMap(t *testing.T) {
	array := []int{1, 2, 3, 4, 5}
	arrNegative := arr.ArrMap[int](array, func(v int) int {
		return 0 - v
	})

	arrExpected := []int{-1, -2, -3, -4, -5}
	assert.Equal(t, arrExpected, arrNegative)

	arrExpected = []int{-1, -2, 0, -4, -5}
	assert.NotEqual(t, arrExpected, arrNegative)
}

func TestArrEach(t *testing.T) {
	array := []int{1, 2, 3, 4, 5}

	sum := 0
	arr.ArrEach[int](array, func(v int) {
		sum += v
	})

	expect := 15
	assert.Equal(t, expect, sum)
}

func TestArrFilter(t *testing.T) {
	array := []int{-1, 3, 0, 4, 2}

	arrFilter := arr.ArrFilter[int](array, func(v int) bool {
		return v > 0
	})
	arrExpected := []int{3, 4, 2}
	assert.Equal(t, arrExpected, arrFilter)

	arrFilter = arr.ArrFilter[int](array, func(v int) bool {
		return v > 4
	})
	arrExpected = []int{}
	assert.Equal(t, arrExpected, arrFilter)
}

func TestArrUnique(t *testing.T) {
	array := []int{2, 3, 0, 3, 2}

	arrUnique := arr.ArrUnique[int](array)
	arrExpected := []int{2, 3, 0}
	assert.Equal(t, arrExpected, arrUnique)
}

func TestArrPrepend(t *testing.T) {
	array := []int{1, 2, 3}

	arrPrepend := arr.ArrPrepend[int](array, 0)
	arrExpected := []int{0, 1, 2, 3}
	assert.Equal(t, arrExpected, arrPrepend)
}

func TestSlice(t *testing.T) {
	array := []int{0, 1, 2, 3, 4, 5}

	arrSlice := arr.ArrSlice[int](array, 2, 4)
	arrExpected := []int{2, 3, 4}
	assert.Equal(t, arrExpected, arrSlice)

	arrSlice = arr.ArrSlice[int](array, 0, 4)
	arrExpected = []int{0, 1, 2, 3, 4}
	assert.Equal(t, arrExpected, arrSlice)

	arrSlice = arr.ArrSlice[int](array, 3, 10)
	arrExpected = []int{3, 4, 5}
	assert.Equal(t, arrExpected, arrSlice)
}

func TestArrAny(t *testing.T) {
	array := []int{0, 1, 2, 3, 4, 5}

	output := arr.ArrAny[int](array)
	assert.IsType(t, []any{}, output)
	assert.Equal(t, output, []any{0, 1, 2, 3, 4, 5})
}

func TestCompare(t *testing.T) {
	arr1 := []int{1, 2, 3}
	arr2 := []int{1, 2, 3}
	arr3 := []int{1, 3, 2}
	assert.Equal(t, true, arr.ArrCompare(arr1, arr2))
	assert.Equal(t, false, arr.ArrCompare(arr1, arr3))
	assert.Equal(t, true, arr.ArrCompareRelaxed(arr1, arr3))
	assert.Equal(t, true, arr.ArrCompareRelaxed(arr2, arr3))
}

func TestIntersect(t *testing.T) {
	arr1 := []int{1, 2, 3, 4, 5, 2, 3, 5, 7, 9, 6, 7}
	arr2 := []int{1, 6, 6, 9, 6, 8, 13}
	arr3 := []int{3, 4, 5, 2, 2, 9, 6}
	expectSliceValues := []int{9, 6}

	intersectResult, _ := arr.ArrIntersect(arr1, arr2, arr3)
	b := arr.ArrCompareRelaxed(intersectResult, expectSliceValues)

	assert.Equal(t, true, b)
}

func TestIntersect2(t *testing.T) {
	arr1 := []int{1, 2, 3}
	arr2 := make([]int, 0)

	intersectResult, _ := arr.ArrIntersect(arr1, arr2)
	b := arr.ArrCompareRelaxed(intersectResult, arr2)

	assert.Equal(t, true, b)
}

func TestIntersect3(t *testing.T) {
	arr1 := []int{1, 2, 3}
	var arr2 []int

	_, err := arr.ArrIntersect(arr1, arr2)
	assert.Equal(t, true, err != nil)
}

func TestIntersectUUID(t *testing.T) {
	arr1 := []string{
		"7f2a68f1-38a9-4625-98d8-053ea8c6aa79",
		"484d836b-b513-464f-8585-8f589c03fb39",
		"f3f4ef4a-266f-4f25-9c25-3102622d835d",
		"f3f4ef4a-266f-4f25-9c25-3102622d835d",
		"f5bc9105-b0d1-483b-904f-929ea9d0b15b",
		"6cf12611-9345-4f98-a426-3c1b841da5e1",
		"78b87535-09cc-4a35-b529-6b8bad2eff4b",
		"78b87535-09cc-4a35-b529-6b8bad2eff4b",
		"a63b246a-fcc2-4cf0-855d-d42a61617ece",
		"f35f1054-a5c8-4723-8d90-45f154222855",
		"484d836b-b513-464f-8585-8f589c03fb39",
		"7f2a68f1-38a9-4625-98d8-053ea8c6aa79",
	}

	arr2 := []string{
		"7f2a68f1-38a9-4625-98d8-053ea8c6aa79",
		"484d836b-b513-464f-8585-8f589c03fb39",
		"f3f4ef4a-266f-4f25-9c25-3102622d835d",
		"f5bc9105-b0d1-483b-904f-929ea9d0b15b",
		"6cf12611-9345-4f98-a426-3c1b841da5e1",
		"78b87535-09cc-4a35-b529-6b8bad2eff4b",
		"a63b246a-fcc2-4cf0-855d-d42a61617ece",
		"f35f1054-a5c8-4723-8d90-45f154222855",
	}

	arr3 := []string{
		"a63b246a-fcc2-4cf0-855d-d42a61617ece",
		"484d836b-b513-464f-8585-8f589c03fb39",
		"484d836b-b513-464f-8585-8f589c03fb30",
	}
	expected := []string{
		"a63b246a-fcc2-4cf0-855d-d42a61617ece",
		"484d836b-b513-464f-8585-8f589c03fb39",
	}

	arrUid1 := arr.ArrMap(arr1, uuid.MustParse)
	arrUid2 := arr.ArrMap(arr2, uuid.MustParse)
	arrUid3 := arr.ArrMap(arr3, uuid.MustParse)
	arrUidExpected := arr.ArrMap(expected, uuid.MustParse)

	b, _ := arr.ArrIntersect(arrUid1, arrUid2, arrUid3)
	ok := arr.ArrCompareRelaxed(b, arrUidExpected)
	assert.Equal(t, true, ok)
}

func TestRemoveDuplicateInt(t *testing.T) {
	intSliceValues := []int{1, 2, 3, 4, 5, 2, 3, 5, 7, 9, 6, 7}
	expectSliceValues := []int{1, 2, 3, 4, 5, 7, 9, 6}

	newSlice := arr.ArrUnique(intSliceValues)

	assert.Equal(t, expectSliceValues, newSlice)
}

func TestRemoveDuplicateUUID(t *testing.T) {
	uuidsStr := []string{
		"7f2a68f1-38a9-4625-98d8-053ea8c6aa79",
		"484d836b-b513-464f-8585-8f589c03fb39",
		"f3f4ef4a-266f-4f25-9c25-3102622d835d",
		"f3f4ef4a-266f-4f25-9c25-3102622d835d",
		"f5bc9105-b0d1-483b-904f-929ea9d0b15b",
		"6cf12611-9345-4f98-a426-3c1b841da5e1",
		"78b87535-09cc-4a35-b529-6b8bad2eff4b",
		"78b87535-09cc-4a35-b529-6b8bad2eff4b",
		"a63b246a-fcc2-4cf0-855d-d42a61617ece",
		"f35f1054-a5c8-4723-8d90-45f154222855",
		"484d836b-b513-464f-8585-8f589c03fb39",
		"7f2a68f1-38a9-4625-98d8-053ea8c6aa79",
	}

	notDuplicate := []string{
		"7f2a68f1-38a9-4625-98d8-053ea8c6aa79",
		"484d836b-b513-464f-8585-8f589c03fb39",
		"f3f4ef4a-266f-4f25-9c25-3102622d835d",
		"f5bc9105-b0d1-483b-904f-929ea9d0b15b",
		"6cf12611-9345-4f98-a426-3c1b841da5e1",
		"78b87535-09cc-4a35-b529-6b8bad2eff4b",
		"a63b246a-fcc2-4cf0-855d-d42a61617ece",
		"f35f1054-a5c8-4723-8d90-45f154222855",
	}
	dupUUID := arr.ArrMap(uuidsStr, uuid.MustParse)
	notDupUUID := arr.ArrMap(notDuplicate, uuid.MustParse)
	removedDup := arr.ArrUnique(dupUUID)
	assert.Equal(t, len(removedDup), len(notDuplicate))
	assert.Equal(t, notDupUUID, removedDup)
}

func TestContainsInt(t *testing.T) {
	intSliceValues := []int{1, 2, 3, 4, 5, 2, 3, 5, 7, 9, 6, 7}
	_, found := arr.ArrFind(intSliceValues, 9)
	assert.Equal(t, true, found)
}

func TestContainsUUID(t *testing.T) {
	uuidsStr := []string{
		"7f2a68f1-38a9-4625-98d8-053ea8c6aa79",
		"484d836b-b513-464f-8585-8f589c03fb39",
		"f3f4ef4a-266f-4f25-9c25-3102622d835d",
		"f5bc9105-b0d1-483b-904f-929ea9d0b15b",
		"6cf12611-9345-4f98-a426-3c1b841da5e1",
		"78b87535-09cc-4a35-b529-6b8bad2eff4b",
		"a63b246a-fcc2-4cf0-855d-d42a61617ece",
		"f35f1054-a5c8-4723-8d90-45f154222855",
	}
	uuidArray := arr.ArrMap(uuidsStr, uuid.MustParse)

	_, found := arr.ArrFind(uuidArray, uuid.MustParse("f35f1054-a5c8-4723-8d90-45f154222855"))
	_, isNotContains := arr.ArrFind(uuidArray, uuid.MustParse("555f1054-a5c8-4723-8d90-45f154222855"))
	assert.Equal(t, true, found)
	assert.Equal(t, false, isNotContains)
}
