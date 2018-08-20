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
	StructDifference
)

type Difference struct {
	LeftVal  interface{}
	RightVal interface{}
}

type Comparison struct {
	Type              DifferenceType
	Index             uint
	Label             string
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
			rightComparator, leftComparator interface{}
			differenceType                  DifferenceType
		)

		comparisons := lVal.Len()
		if rVal.Len() > comparisons {
			comparisons = rVal.Len()
		}

		for i := 0; i < comparisons; i++ {
			differenceType = SliceDifference

			if i < rVal.Len() {
				rightComparator = rVal.Index(i).Interface()
			} else {
				rightComparator = nil
				differenceType = SliceAdditionalValue
			}

			if i < lVal.Len() {
				leftComparator = lVal.Index(i).Interface()
			} else {
				leftComparator = nil
				differenceType = SliceAdditionalValue
			}

			comparison := compare(leftComparator, rightComparator)
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
	case reflect.Struct:
		leftV := reflect.ValueOf(left)
		for i := 0; i < leftV.NumField(); i++ {
			leftName := leftV.Type().Field(i).Name
			leftComparitor := leftV.Field(i).Interface()
			rightComparitor := reflect.ValueOf(right).FieldByName(leftName).Interface()
			comparison := compare(leftComparitor, rightComparitor)
			if !reflect.DeepEqual(comparison, Comparison{}) {
				comparison.Type = StructDifference
				comparison.Label = leftName
				diffs.DifferenceDetails = append(diffs.DifferenceDetails, comparison)
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
