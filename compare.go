package main

import (
	"reflect"
)

type DifferenceType int

const (
	RawDifference DifferenceType = iota
	TypeDifference
	SliceDifference
	SliceAdditionalValue
)

type Difference struct {
	LeftVal  interface{}
	RightVal interface{}
}

type Comparison struct {
	Type              DifferenceType
	Index             uint
	Difference        Difference
	DifferenceDetails []Comparison
}

func Compare(left, right interface{}) Comparison {
	return compare(left, right)
}

func compare(left, right interface{}) (diffs Comparison) {
	if reflect.DeepEqual(left, right) {
		return
	}

	lVal := reflect.ValueOf(left)
	rVal := reflect.ValueOf(right)

	lKind := lVal.Kind()

	if lKind != rVal.Kind() {
		diffs.Type = TypeDifference
		diffs.Difference = Difference{
			LeftVal:  left,
			RightVal: right,
		}

		return
	}

	switch lKind {
	case reflect.Slice:
		var (
			rightComparator interface{}
			differenceType DifferenceType
		)

		comparisons := lVal.Len()
		if rVal.Len() > comparisons {
			comparisons = rVal.Len()
		}

		for i := 0; i < comparisons; i++ {
			if i < rVal.Len() {
				rightComparator = rVal.Index(i).Interface()
				differenceType = SliceDifference
			} else {
				rightComparator = nil
				differenceType = SliceAdditionalValue
			}

			comparison := compare(lVal.Index(i).Interface(), rightComparator)
			if !reflect.DeepEqual(comparison, Comparison{}) {
				comparison.Type = differenceType
				comparison.Index = uint(i)
				diffs.DifferenceDetails = append(diffs.DifferenceDetails, comparison)
			}
		}
	case reflect.Bool:
		if lVal != rVal {
			diffs.Difference = Difference{
				LeftVal:  lVal.Bool(),
				RightVal: rVal.Bool(),
			}
		}
	case reflect.Int:
		if left != right {
			diffs.Difference = Difference{
				LeftVal:  left,
				RightVal: right,
			}
		}
	case reflect.String:
		if left != right {
			diffs.Difference = Difference{
				LeftVal:  left,
				RightVal: right,
			}
		}
	default:
		//diffs.Difference = Difference{
		//	LeftVal:  left,
		//	RightVal: right,
		//}
	}

	return
}
