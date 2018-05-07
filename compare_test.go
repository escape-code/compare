package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	. "github.com/simonjjones/compare"
)

var _ = Describe("Compare", func() {
	DescribeTable("basic types",
		func(l, r interface{}, expectedComparison Comparison) {
			c := Compare(l, r)
			Expect(c).To(Equal(expectedComparison))
		},
		Entry("returns empty comparison when comparing two equal bools", true, true, nil),
		Entry("returns comparison with one difference when comparing two unequal bools",
			true,
			false,
			Comparison{
				Difference: Difference{
					LeftVal:  true,
					RightVal: false,
				},
			},
		),
		Entry("returns empty comparison comparing two equal ints", 2345, 2345, nil),
		Entry("returns comparison with one difference when comparing two unequal ints",
			2345,
			6789,
			Comparison{
				Difference: Difference{
					LeftVal:  2345,
					RightVal: 6789,
				},
			},
		),
		Entry("returns empty comparison comparing two equal strings", "some writing", "some writing", nil),
		Entry("returns comparison with one difference when comparing two unequal strings",
			"some writing",
			"about a thing",
			Comparison{
				Difference: Difference{
					LeftVal:  "some writing",
					RightVal: "about a thing",
				},
			},
		),
		Entry("returns empty comparison when comparing nils", nil, nil, nil),
		Entry("returns comparison with one difference when comparing nil with bool",
			nil,
			true,
			Comparison{
				Type: TypeDifference,
				Difference: Difference{
					LeftVal:  nil,
					RightVal: true,
				},
			},
		),
		Entry("returns comparison with one difference when comparing bool with nil",
			true,
			nil,
			Comparison{
				Type: TypeDifference,
				Difference: Difference{
					LeftVal:  true,
					RightVal: nil,
				},
			},
		),
	)

	Describe("slices", func() {
		It("returns empty comparison when comparing two equal slices", func() {
			leftSlice := []bool{true, false, true}
			rightSlice := []bool{true, false, true}

			c := Compare(leftSlice, rightSlice)
			Expect(c).To(Equal(Comparison{}))
		})

		It("returns comparison with difference details when one value is different between the slices", func() {
			leftSlice := []bool{true, false, true}
			rightSlice := []bool{true, false, false}

			c := Compare(leftSlice, rightSlice)
			Expect(c).To(Equal(Comparison{
				DifferenceDetails: []Comparison{
					{
						Type:  SliceDifference,
						Index: uint(2),
						Difference: Difference{
							LeftVal:  true,
							RightVal: false,
						},
					},
				},
			}))
		})

		It("returns comparison with several difference details when several values are different between the slices", func() {
			leftSlice := []bool{false, false, true}
			rightSlice := []bool{true, false, false}

			c := Compare(leftSlice, rightSlice)
			Expect(c).To(Equal(Comparison{
				DifferenceDetails: []Comparison{
					{
						Type:  SliceDifference,
						Index: uint(0),
						Difference: Difference{
							LeftVal:  false,
							RightVal: true,
						},
					},
					{
						Type:  SliceDifference,
						Index: uint(2),
						Difference: Difference{
							LeftVal:  true,
							RightVal: false,
						},
					},
				},
			}))
		})

		It("returns comparison with addition values as differences when left slice is larger than right", func() {
			leftSlice := []int{1, 2, 3, 4}
			rightSlice := []int{1, 2}

			c := Compare(leftSlice, rightSlice)
			Expect(c).To(Equal(Comparison{
				DifferenceDetails: []Comparison{
					{
						Type: SliceAdditionalValue,
						Index: uint(2),
						Difference: Difference{
							LeftVal: 3,
							RightVal: nil,
						},
					},
					{
						Type: SliceAdditionalValue,
						Index: uint(3),
						Difference: Difference{
							LeftVal: 4,
							RightVal: nil,
						},
					},
				},
			}))
		})

		It("returns comparison with addition values as differences when right slice is larger than left", func() {
			leftSlice := []string{"a", "b"}
			rightSlice := []string{"a", "b", "c", "d"}

			c := Compare(leftSlice, rightSlice)
			Expect(c).To(Equal(Comparison{
				DifferenceDetails: []Comparison{
					{
						Type: SliceAdditionalValue,
						Index: uint(2),
						Difference: Difference{
							LeftVal: nil,
							RightVal: "c",
						},
					},
					{
						Type: SliceAdditionalValue,
						Index: uint(3),
						Difference: Difference{
							LeftVal: nil,
							RightVal: "d",
						},
					},
				},
			}))
		})
	})
})
