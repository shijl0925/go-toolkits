package slice_test

import (
	"github.com/shijl0925/go-toolkits/slice"
	"reflect"
	"testing"
)

func Test_SumSlice(t *testing.T) {
	var tests = []struct {
		name     string
		slice    []float64
		expected float64
	}{
		{"test1", []float64{1, 2, 3, 4, 5}, 15},
		{"test2", []float64{1.1, 2.0, 3.5}, 6.6},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := slice.Sum(tc.slice); got != tc.expected {
				t.Errorf("Sum() expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func Test_MaxSlice(t *testing.T) {
	var tests = []struct {
		name     string
		slice    []float64
		expected float64
	}{
		{"test1", []float64{-1, -5, -3, -4, -2}, -1},
		{"test2", []float64{1.1, 3.5, 2.0, 0.1}, 3.5},
		{"test3", []float64{-1, 5, 3, 4, 2}, 5},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := slice.Max(tc.slice); got != tc.expected {
				t.Errorf("Max() expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func Test_MinSlice(t *testing.T) {
	var tests = []struct {
		name     string
		slice    []float64
		expected float64
	}{
		{"test1", []float64{-1, 5, 3, 4, 2}, -1},
		{"test2", []float64{1.1, 3.5, 2.0, 0.1}, 0.1},
		{"test3", []float64{-1, -5, -3, -4, -2}, -5},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := slice.Min(tc.slice); got != tc.expected {
				t.Errorf("Min() expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func Test_Insert(t *testing.T) {
	var tests = []struct {
		name     string
		slice    []int
		element  int
		index    int
		expected []int
	}{
		{"test1", []int{}, 1, 0, []int{1}},
		{"test2", []int{1, 2, 3, 4}, 2, 1, []int{1, 2, 2, 3, 4}},
		{"test3", []int{1, 2, 3, 4, 5, 6}, 7, 6, []int{1, 2, 3, 4, 5, 6, 7}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := slice.Insert(tc.slice, tc.element, tc.index); !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("Insert() expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func Test_Add(t *testing.T) {
	var tests = []struct {
		name     string
		slice    []int
		element  int
		index    int
		expected []int
	}{
		{"test1", []int{}, 1, 0, []int{1}},
		{"test2", []int{1, 2, 3, 4}, 2, 1, []int{1, 2, 2, 3, 4}},
		{"test3", []int{1, 2, 3, 4, 5, 6}, 7, 6, []int{1, 2, 3, 4, 5, 6, 7}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := slice.AddV2(tc.slice, tc.element, tc.index); !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("Add() expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func Test_InsertMany(t *testing.T) {
	var tests = []struct {
		name     string
		slice    []int
		elements []int
		index    int
		expected []int
	}{
		{"test1", []int{}, []int{1}, 0, []int{1}},
		{"test2", []int{1, 2, 3, 4}, []int{2, 3}, 1, []int{1, 2, 3, 2, 3, 4}},
		{"test3", []int{1, 2, 3, 4, 5, 6}, []int{7}, 6, []int{1, 2, 3, 4, 5, 6, 7}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := slice.InsertMany(tc.slice, tc.elements, tc.index); !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("InsertMany() expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func Test_AddMany(t *testing.T) {
	var tests = []struct {
		name     string
		slice    []int
		elements []int
		index    int
		expected []int
	}{
		{"test1", []int{}, []int{1}, 0, []int{1}},
		{"test2", []int{1, 2, 3, 4}, []int{2, 3}, 1, []int{1, 2, 3, 2, 3, 4}},
		{"test3", []int{1, 2, 3, 4, 5, 6}, []int{7}, 6, []int{1, 2, 3, 4, 5, 6, 7}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := slice.AddMany(tc.slice, tc.elements, tc.index); !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("AddMany() expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func Test_Pop(t *testing.T) {
	var s, want []string
	var lastWant string
	t.Run("test0", func(t *testing.T) {
		if last, got := slice.Pop(s); !reflect.DeepEqual(got, want) {
			t.Errorf("Pop() expected %v, got %v", want, got)
			if last != lastWant {
				t.Errorf("Pop() element expected %v, got %v", lastWant, last)
			}
		}

	})
	var tests = []struct {
		name     string
		slice    []int
		expected []int
	}{
		{"test1", []int{1, 2, 3, 4}, []int{1, 2, 3}},
		{"test2", []int{1, 2, 3, 4, 5}, []int{1, 2, 3, 4}},
		{"test3", []int{1, 2, 3, 4, 5, 6}, []int{1, 2, 3, 4, 5}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if last, got := slice.Pop(tc.slice); !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("Pop() expected %v, got %v", tc.expected, got)
				if last != tc.slice[len(tc.slice)-1] {
					t.Errorf("Pop() element expected %v, got %v", last, tc.slice[len(tc.slice)-1])
				}
			}
		})
	}
}

func Test_Delete(t *testing.T) {
	var tests = []struct {
		name     string
		slice    []int
		index    int
		expected []int
	}{
		{"test1", []int{1, 2, 3, 4}, 1, []int{1, 3, 4}},
		{"test2", []int{1, 2, 3, 4}, 3, []int{1, 2, 3}},
		{"test4", []int{1, 2, 3, 4, 5, 6}, 2, []int{1, 2, 4, 5, 6}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := slice.Delete(tc.slice, tc.index); !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("Delete() expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func Test_DeleteMany(t *testing.T) {
	var tests = []struct {
		name     string
		slice    []int
		start    int
		end      int
		expected []int
	}{
		{"test1", []int{1, 2, 3, 4}, 1, 2, []int{1, 3, 4}},
		{"test2", []int{1, 2, 3, 4, 5}, 2, 4, []int{1, 2, 5}},
		{"test3", []int{1, 2, 3, 4, 5, 6}, 4, 5, []int{1, 2, 3, 4, 6}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := slice.DeleteMany(tc.slice, tc.start, tc.end); !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("DeleteMany() expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func Test_Append(t *testing.T) {
	var tests = []struct {
		name     string
		slice    []int
		element  int
		expected []int
	}{
		{"test1", []int{}, 1, []int{1}},
		{"test2", []int{1, 2, 3, 4}, 5, []int{1, 2, 3, 4, 5}},
		{"test3", []int{1, 2, 3, 4, 5, 6}, 7, []int{1, 2, 3, 4, 5, 6, 7}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := slice.Append(tc.slice, tc.element); !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("Append() expected %v, got %v", tc.expected, got)
			}
		})
	}

	slice1 := []int{1, 2, 3, 4}
	elements := []int{5, 6, 7}
	want := []int{1, 2, 3, 4, 5, 6, 7}
	t.Run("test0", func(t *testing.T) {
		if got := slice.Append(slice1, elements...); !reflect.DeepEqual(got, want) {
			t.Errorf("Append() expected %v, got %v", want, got)
		}
	})
}

func Test_Extend(t *testing.T) {
	var tests = []struct {
		name     string
		slice    []int
		elements []int
		expected []int
	}{
		{"test1", []int{}, []int{5, 6, 7}, []int{5, 6, 7}},
		{"test2", []int{1, 2, 3, 4}, []int{5, 6, 7}, []int{1, 2, 3, 4, 5, 6, 7}},
		{"test3", []int{1, 2, 3, 4, 5, 6}, []int{7}, []int{1, 2, 3, 4, 5, 6, 7}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := slice.Extend(tc.slice, tc.elements); !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("Extend() expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func Test_Unique(t *testing.T) {
	var tests = []struct {
		name     string
		slice    []int
		expected []int
	}{
		{"test1", []int{}, []int{}},
		{"test2", []int{1, 2, 3, 4}, []int{1, 2, 3, 4}},
		{"test3", []int{1, 2, 2, 3, 4, 4, 5}, []int{1, 2, 3, 4, 5}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := slice.Unique(tc.slice); !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("Unique() expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func Test_ReverseSelfSlice(t *testing.T) {
	s1 := []int{1, 2, 3, 4, 5}
	t.Run("test1", func(t *testing.T) {
		slice.ReverseSelf(s1)
		if !reflect.DeepEqual(s1, []int{5, 4, 3, 2, 1}) {
			t.Errorf("Reverse() expected %v, got %v", 0, s1)
		}
	})

	s2 := []string{"one", "two", "three"}
	t.Run("test2", func(t *testing.T) {
		slice.ReverseSelf(s2)
		if !reflect.DeepEqual(s2, []string{"three", "two", "one"}) {
			t.Errorf("Reverse() expected %v, got %v", []string{"three", "two", "one"}, s2)
		}
	})

	s3 := []byte("Google")
	t.Run("test3", func(t *testing.T) {
		slice.ReverseSelf(s3)
		if string(s3) != "elgooG" {
			t.Errorf("Reverse() expected %v, got %v", "elgooG", string(s3))
		}
	})
}
func Test_ReverseSlice(t *testing.T) {
	t.Run("test1", func(t *testing.T) {
		if got := slice.Reverse([]int{1, 2, 3, 4, 5}); !reflect.DeepEqual(got, []int{5, 4, 3, 2, 1}) {
			t.Errorf("Reverse() expected %v, got %v", []int{5, 4, 3, 2, 1}, got)
		}
	})
	t.Run("test2", func(t *testing.T) {
		if got := slice.Reverse([]string{"one", "two", "three"}); !reflect.DeepEqual(got, []string{"three", "two", "one"}) {
			t.Errorf("Reverse() expected %v, got %v", []string{"three", "two", "one"}, got)
		}
	})
	t.Run("test3", func(t *testing.T) {
		if got := slice.Reverse([]byte("Google")); string(got) != "elgooG" {
			t.Errorf("Reverse() expected %v, got %v", "elgooG", string(got))
		}
	})
}

func Test_RepeatSlice(t *testing.T) {
	t.Run("test1", func(t *testing.T) {
		if got := slice.Repeat(1, 5); !reflect.DeepEqual(got, []int{1, 1, 1, 1, 1}) {
			t.Errorf("Repeat() expected %v, got %v", []int{1, 1, 1, 1, 1}, got)
		}
	})
	t.Run("test2", func(t *testing.T) {
		wanted := [][]int{[]int{1, 2, 3}, []int{1, 2, 3}, []int{1, 2, 3}, []int{1, 2, 3}, []int{1, 2, 3}}
		if got := slice.Repeat([]int{1, 2, 3}, 5); !reflect.DeepEqual(got, wanted) {
			t.Errorf("Repeat() expected %v, got %v", wanted, got)
		}
	})
	t.Run("test3", func(t *testing.T) {
		wanted := []string{"one", "one", "one"}
		if got := slice.Repeat("one", 3); !reflect.DeepEqual(got, wanted) {
			t.Errorf("Repeat() expected %v, got %v", wanted, got)
		}
	})
}

func Test_Product(t *testing.T) {
	tests := []struct {
		name   string
		slices [][]string
		want   [][]string
	}{
		{
			name:   "test1",
			slices: [][]string{[]string{"A", "B", "C"}, []string{"one", "two", "three"}},
			want: [][]string{
				{"A", "one"},
				{"A", "two"},
				{"A", "three"},
				{"B", "one"},
				{"B", "two"},
				{"B", "three"},
				{"C", "one"},
				{"C", "two"},
				{"C", "three"},
			},
		},
		{
			name:   "test2",
			slices: [][]string{[]string{"A", "B"}, []string{"C", "D"}, []string{"E", "F"}},
			want: [][]string{
				{"A", "C", "E"},
				{"A", "C", "F"},
				{"A", "D", "E"},
				{"A", "D", "F"},
				{"B", "C", "E"},
				{"B", "C", "F"},
				{"B", "D", "E"},
				{"B", "D", "F"},
			},
		},
		{
			name:   "test3",
			slices: [][]string{[]string{"A", "B"}, []string{"C", "D", "E", "F"}},
			want: [][]string{
				{"A", "C"},
				{"A", "D"},
				{"A", "E"},
				{"A", "F"},
				{"B", "C"},
				{"B", "D"},
				{"B", "E"},
				{"B", "F"},
			},
		},
		{
			name:   "test4",
			slices: [][]string{[]string{"A", "B", "C", "D", "E"}, []string{"F", "G"}},
			want: [][]string{
				{"A", "F"},
				{"A", "G"},
				{"B", "F"},
				{"B", "G"},
				{"C", "F"},
				{"C", "G"},
				{"D", "F"},
				{"D", "G"},
				{"E", "F"},
				{"E", "G"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := slice.Product(tt.slices...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Product() = %v, want %v", got, tt.want)
			}
		})
	}
	t.Run("test5", func(t *testing.T) {
		slices := [][]int{[]int{1, 2, 3}}
		wanted := [][]int{{1}, {2}, {3}}
		if got := slice.Product(slices...); !reflect.DeepEqual(got, wanted) {
			t.Errorf("Product() expected %v, got %v", wanted, got)
		}
	})
	t.Run("test6", func(t *testing.T) {
		var slices [][]int
		var wanted [][]int
		if got := slice.Product(slices...); !reflect.DeepEqual(got, wanted) {
			t.Errorf("Product() expected %v, got %v", wanted, got)
		}
	})
}

func Test_DiffSet(t *testing.T) {
	var tests = []struct {
		name     string
		slice1   []int
		slice2   []int
		expected []int
	}{
		{"test1", []int{1, 2, 3, 4}, []int{3, 5, 4, 6, 7}, []int{1, 2}},
		{"test2", []int{1, 3, 4, 2}, []int{}, []int{1, 3, 4, 2}},
		{"test3", []int{1, 3, 2, 4}, []int{5, 6, 7, 8}, []int{1, 3, 2, 4}},
		{"test4", []int{}, []int{3, 4, 5, 6, 7}, []int{}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := slice.DiffSet(tc.slice1, tc.slice2); !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("DiffSet() expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func Test_IntersectSet(t *testing.T) {
	var tests = []struct {
		name     string
		slice1   []int
		slice2   []int
		expected []int
	}{
		{"test1", []int{1, 2, 3, 4}, []int{3, 5, 4, 6, 7}, []int{3, 4}},
		{"test2", []int{1, 3, 4, 2}, []int{}, []int{}},
		{"test3", []int{1, 3, 2, 4}, []int{5, 6, 7, 8}, []int{}},
		{"test4", []int{}, []int{3, 4, 5, 6, 7}, []int{}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := slice.IntersectSet(tc.slice1, tc.slice2); !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("IntersectSet() expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func Test_UnionSet(t *testing.T) {
	var tests = []struct {
		name     string
		slice1   []int
		slice2   []int
		expected []int
	}{
		{"test1", []int{1, 2, 3, 4}, []int{3, 5, 4, 6, 7}, []int{1, 2, 3, 4, 5, 6, 7}},
		{"test2", []int{1, 3, 4, 2}, []int{}, []int{1, 3, 4, 2}},
		{"test3", []int{1, 3, 2, 4}, []int{5, 6, 7, 8}, []int{1, 3, 2, 4, 5, 6, 7, 8}},
		{"test4", []int{}, []int{3, 4, 5, 6, 7}, []int{3, 4, 5, 6, 7}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := slice.UnionSet(tc.slice1, tc.slice2); !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("UnionSet() expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func Test_Index(t *testing.T) {
	var tests = []struct {
		name     string
		slice    []int
		e        int
		expected int
	}{
		{"test1", []int{1, 2, 3, 4}, 1, 0},
		{"test2", []int{1, 2, 3, 4}, -3, -1},
		{"test3", []int{1, 2, 3, 4, 4}, 4, 3},
		{"test4", []int{1, 2, 3, 4, 5, 4, 6}, 4, 3},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := slice.Index(tc.slice, tc.e); !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("Index() expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func Test_LastIndex(t *testing.T) {
	var tests = []struct {
		name     string
		slice    []int
		e        int
		expected int
	}{
		{"test1", []int{1, 2, 3, 4}, 2, 1},
		{"test2", []int{1, 2, 3, 4}, -3, -1},
		{"test3", []int{1, 2, 3, 4, 4}, 4, 4},
		{"test4", []int{1, 2, 3, 4, 5, 4, 6}, 4, 5},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := slice.LastIndex(tc.slice, tc.e); !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("LastIndex() expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func Test_IndexAll(t *testing.T) {
	var tests = []struct {
		name     string
		slice    []int
		e        int
		expected []int
	}{
		{"test1", []int{1, 2, 3, 4}, 2, []int{1}},
		{"test2", []int{1, 2, 3, 4}, -3, []int{}},
		{"test3", []int{1, 2, 3, 4, 4}, 4, []int{3, 4}},
		{"test4", []int{1, 2, 3, 4, 5, 4, 6}, 4, []int{3, 5}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := slice.IndexAll(tc.slice, tc.e); !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("IndexAll() expected %v, got %v", tc.expected, got)
			}
		})
	}
}
